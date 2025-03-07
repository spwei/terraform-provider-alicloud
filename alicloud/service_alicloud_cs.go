package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"encoding/base64"

	"encoding/json"

	"github.com/alibabacloud-go/cs-20151215/v3/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
)

type CsService struct {
	client *connectivity.AliyunClient
}

type CsClient struct {
	client *client.Client
}

type Component struct {
	ComponentName string `json:"component_name"`
	Version       string `json:"version"`
	NextVersion   string `json:"next_version"`
	CanUpgrade    bool   `json:"can_upgrade"`
	Required      bool   `json:"required"`
	Status        string `json:"status"`
	ErrMessage    string `json:"err_message"`
	Config        string `json:"config"`
	ConfigSchema  string `json:"config_schema"`
}

const (
	COMPONENT_AUTO_SCALER      = "cluster-autoscaler"
	COMPONENT_DEFAULT_VRESION  = "v1.0.0"
	SCALING_CONFIGURATION_NAME = "kubernetes_autoscaler_autogen"
	DefaultECSTag              = "k8s.aliyun.com"
	DefaultClusterTag          = "ack.aliyun.com"
	RECYCLE_MODE_LABEL         = "k8s.io/cluster-autoscaler/node-template/label/policy"
	DefaultAutoscalerTag       = "k8s.io/cluster-autoscaler"
	SCALING_GROUP_NAME         = "sg-%s-%s"
	DEFAULT_COOL_DOWN_TIME     = 300
	RELEASE_MODE               = "release"
	RECYCLE_MODE               = "recycle"

	PRIORITY_POLICY       = "PRIORITY"
	COST_OPTIMIZED_POLICY = "COST_OPTIMIZED"
	BALANCE_POLICY        = "BALANCE"

	UpgradeClusterTimeout = 30 * time.Minute

	IdMsgWithTask = IdMsg + "TaskInfo: %s" // wait for async task info
)

var (
	ATTACH_SCRIPT_WITH_VERSION = `#!/bin/sh
curl http://aliacs-k8s-%s.oss-%s.aliyuncs.com/public/pkg/run/attach/%s/attach_node.sh | bash -s -- --openapi-token %s --ess true `
	NETWORK_ADDON_NAMES = []string{"terway", "kube-flannel-ds", "terway-eni", "terway-eniip"}
)

func (s *CsService) GetContainerClusterByName(name string) (cluster cs.ClusterType, err error) {
	name = Trim(name)
	invoker := NewInvoker()
	var clusters []cs.ClusterType
	err = invoker.Run(func() error {
		raw, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeClusters(name)
		})
		if e != nil {
			return e
		}
		clusters, _ = raw.([]cs.ClusterType)
		return nil
	})

	if err != nil {
		return cluster, fmt.Errorf("Describe cluster failed by name %s: %#v.", name, err)
	}

	if len(clusters) < 1 {
		return cluster, GetNotFoundErrorFromString(GetNotFoundMessage("Container Cluster", name))
	}

	for _, c := range clusters {
		if c.Name == name {
			return c, nil
		}
	}
	return cluster, GetNotFoundErrorFromString(GetNotFoundMessage("Container Cluster", name))
}

func (s *CsService) GetContainerClusterAndCertsByName(name string) (*cs.ClusterType, *cs.ClusterCerts, error) {
	cluster, err := s.GetContainerClusterByName(name)
	if err != nil {
		return nil, nil, err
	}
	var certs cs.ClusterCerts
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.GetClusterCerts(cluster.ClusterID)
		})
		if e != nil {
			return e
		}
		certs, _ = raw.(cs.ClusterCerts)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &cluster, &certs, nil
}

func (s *CsService) DescribeContainerApplication(clusterName, appName string) (app cs.GetProjectResponse, err error) {
	appName = Trim(appName)
	cluster, certs, err := s.GetContainerClusterAndCertsByName(clusterName)
	if err != nil {
		return app, err
	}
	raw, err := s.client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
		return csProjectClient.GetProject(appName)
	})
	app, _ = raw.(cs.GetProjectResponse)
	if err != nil {
		if IsExpectedErrors(err, []string{"Not Found"}) {
			return app, GetNotFoundErrorFromString(GetNotFoundMessage("Container Application", appName))
		}
		return app, fmt.Errorf("Getting Application failed by name %s: %#v.", appName, err)
	}
	if app.Name != appName {
		return app, GetNotFoundErrorFromString(GetNotFoundMessage("Container Application", appName))
	}
	return
}

func (s *CsService) WaitForContainerApplication(clusterName, appName string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		app, err := s.DescribeContainerApplication(clusterName, appName)
		if err != nil {
			return err
		}

		if strings.ToLower(app.CurrentState) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for container application %s is timeout and current status is %s.", string(status), app.CurrentState))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *CsService) DescribeCsKubernetes(id string) (cluster *cs.KubernetesClusterDetail, err error) {
	invoker := NewInvoker()
	var requestInfo *cs.Client
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeKubernetesClusterDetail(id)
		})
		response = raw
		return err
	}); err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, WrapErrorf(err, NotFoundMsg, DenverdinoAliyungo)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, "DescribeKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = id
		addDebug("DescribeKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ = response.(*cs.KubernetesClusterDetail)
	if cluster.ClusterId != id {
		return cluster, WrapErrorf(Error(GetNotFoundMessage("CsKubernetes", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

// DescribeClusterKubeConfig return cluster kube_config credential.
// It's used for kubernetes/managed_kubernetes/serverless_kubernetes.
// Deprecated, use CsClient.DescribeClusterKubeConfigWithExpiration
func (s *CsService) DescribeClusterKubeConfig(clusterId string, isResource bool) (*cs.ClusterConfig, error) {
	invoker := NewInvoker()
	var response interface{}
	var requestInfo *cs.Client
	var config *cs.ClusterConfig
	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeClusterUserConfig(clusterId, false)
		})
		response = raw
		return err
	}); err != nil {
		if isResource {
			return nil, WrapErrorf(err, DefaultErrorMsg, clusterId, "DescribeClusterUserConfig", DenverdinoAliyungo)
		}
		return nil, WrapErrorf(err, DataDefaultErrorMsg, clusterId, "DescribeClusterUserConfig", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = clusterId
		addDebug("DescribeClusterUserConfig", response, requestInfo, requestMap)
	}
	config, _ = response.(*cs.ClusterConfig)
	return config, nil
}

// DescribeClusterKubeConfigWithExpiration return cluster kube_config credential with expiration time.
// It's used for kubernetes/managed_kubernetes/serverless_kubernetes.
func (s *CsClient) DescribeClusterKubeConfigWithExpiration(clusterId string, temporaryDurationMinutes int64) (*client.DescribeClusterUserKubeconfigResponseBody, error) {
	if clusterId == "" {
		return nil, WrapError(fmt.Errorf("clusterid is empty"))
	}

	request := &client.DescribeClusterUserKubeconfigRequest{
		PrivateIpAddress: tea.Bool(false),
	}
	if temporaryDurationMinutes > 0 {
		request.TemporaryDurationMinutes = tea.Int64(temporaryDurationMinutes)
	}
	kubeConfig, err := s.client.DescribeClusterUserKubeconfig(tea.String(clusterId), request)
	if err != nil {
		return nil, WrapError(err)
	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = clusterId
		addDebug("DescribeClusterUserConfig", kubeConfig, request, requestMap)
	}

	return kubeConfig.Body, nil
}

// This function returns all available addons metadata of the cluster
func (s *CsClient) DescribeClusterAddonsMetadata(clusterId string) (map[string]*Component, error) {
	result := make(map[string]*Component)

	resp, err := s.client.DescribeClusterAddonsVersion(&clusterId)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeClusterAddonsVersion", err)
	}

	for name, addon := range resp.Body {
		fields := addon.(map[string]interface{})
		version := fields["version"].(string)
		nextVersion := fields["next_version"].(string)
		canUpgrade := fields["can_upgrade"].(bool)
		required := fields["required"].(bool)
		config := fields["config"].(string)
		c := &Component{
			ComponentName: name,
			Version:       version,
			NextVersion:   nextVersion,
			CanUpgrade:    canUpgrade,
			Required:      required,
			Config:        config,
		}
		result[name] = c
	}

	return result, nil
}

// This function returns the latest addon status information
func (s *CsClient) DescribeCsKubernetesAddonStatus(clusterId string, addonName string) (*Component, error) {
	result := &Component{}

	resp, err := s.client.DescribeClusterAddonsUpgradeStatus(&clusterId, &client.DescribeClusterAddonsUpgradeStatusRequest{
		ComponentIds: []*string{tea.String(addonName)},
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeClusterAddonsUpgradeStatus", err)
	}

	addon, ok := resp.Body[addonName]
	if !ok {
		return nil, WrapErrorf(Error(GetNotFoundMessage("alicloud_cs_kubernetes_addon", addonName)), ResourceNotfound)
	}

	addonInfo := addon.(map[string]interface{})["addon_info"]
	tasks := addon.(map[string]interface{})["tasks"]
	result.Version = addonInfo.(map[string]interface{})["version"].(string)
	result.CanUpgrade = addon.(map[string]interface{})["can_upgrade"].(bool)
	result.Status = tasks.(map[string]interface{})["status"].(string)
	if message, ok := tasks.(map[string]interface{})["message"]; ok {
		result.ErrMessage = message.(string)
	}

	if result.Version == "" {
		return result, WrapErrorf(Error(GetNotFoundMessage("alicloud_cs_kubernetes_addon", addonName)), ResourceNotfound)
	}

	return result, nil
}

// This function returns the latest addon instance
func (s *CsClient) DescribeCsKubernetesAddonInstance(clusterId string, addonName string) (*Component, error) {
	component := &Component{}

	resp, err := s.client.DescribeClusterAddonInstance(&clusterId, tea.String(addonName))
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddonInstance", err)
	}

	component.Config = *resp.Body.Config

	return component, nil
}

// This function returns the status of all available addons of the cluster
func (s *CsClient) DescribeCsKubernetesAllAvailableAddons(clusterId string) (map[string]*Component, error) {
	availableAddons, err := s.DescribeClusterAddonsMetadata(clusterId)
	if err != nil {
		return nil, err
	}

	queryList := make([]*string, 0)
	for name := range availableAddons {
		queryList = append(queryList, tea.String(name))
	}

	status, err := s.DescribeCsKubernetesAllAddonsStatus(clusterId, queryList)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesExistedAddons", err)
	}

	addonInstances := make(map[string]*Component)
	for name := range availableAddons {
		addonInstance, err := s.DescribeCsKubernetesAddonInstance(clusterId, name)
		if err != nil {
			if e, ok := err.(*ComplexError); ok {
				if sdkError, ok := e.Cause.(*tea.SDKError); ok && regexp.MustCompile(NotFound).MatchString(tea.StringValue(sdkError.Code)) {
					log.Printf("[DEBUG] %s addon instance %s not found.", clusterId, name)
					continue
				}
			}
			return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesExistedAddons", err)
		}
		addonInstances[name] = addonInstance
	}

	for name, addon := range availableAddons {
		if _, ok := status[name]; !ok {
			continue
		}
		addon.Version = status[name].Version
		addon.CanUpgrade = status[name].CanUpgrade
		addon.ErrMessage = status[name].ErrMessage
		if _, ok := addonInstances[name]; ok {
			addon.Config = addonInstances[name].Config
			addon.Status = addonInstances[name].Status
		} else {
			addon.Config = status[name].Config
			addon.Status = status[name].Status
		}
	}

	return availableAddons, nil
}

// This function returns the latest multiple addons status information
func (s *CsClient) DescribeCsKubernetesAllAddonsStatus(clusterId string, addons []*string) (map[string]*Component, error) {
	addonsStatus := make(map[string]*Component)
	resp, err := s.client.DescribeClusterAddonsUpgradeStatus(&clusterId, &client.DescribeClusterAddonsUpgradeStatusRequest{
		ComponentIds: addons,
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeClusterAddonsUpgradeStatus", err)
	}

	for name, status := range resp.Body {
		c := &Component{}
		addonInfo := status.(map[string]interface{})["addon_info"]
		tasks := status.(map[string]interface{})["tasks"]
		c.Version = addonInfo.(map[string]interface{})["version"].(string)
		c.CanUpgrade = status.(map[string]interface{})["can_upgrade"].(bool)
		c.Status = tasks.(map[string]interface{})["status"].(string)
		if message, ok := tasks.(map[string]interface{})["message"]; ok {
			c.ErrMessage = message.(string)
		}

		addonsStatus[name] = c
	}

	return addonsStatus, nil
}

func (s *CsClient) DescribeCsKubernetesAddon(id string) (*Component, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	clusterId := parts[0]
	addonName := parts[1]
	addonsMetadata, err := s.DescribeClusterAddonsMetadata(clusterId)
	if err != nil {
		return nil, err
	}

	addonStatus, err := s.DescribeCsKubernetesAddonStatus(clusterId, addonName)
	if err != nil {
		return nil, err
	}

	addonInstance, err := s.DescribeCsKubernetesAddonInstance(clusterId, addonName)
	if err != nil {
		return nil, err
	}

	// Update some fields
	if addon, existed := addonsMetadata[addonName]; existed {
		addon.Version = addonStatus.Version
		addon.Status = addonStatus.Status
		addon.CanUpgrade = addonStatus.CanUpgrade
		addon.ErrMessage = addonStatus.ErrMessage
		addon.Config = addonInstance.Config
		return addon, nil
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("alicloud_cs_kubernetes_addon", id)), ResourceNotfound)
}

func (s *CsClient) CsKubernetesAddonStateRefreshFunc(clusterId string, addonName string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetesAddonStatus(clusterId, addonName)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatusWithResponse, clusterId, object.ErrMessage))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CsClient) CsKubernetesAddonExistRefreshFunc(clusterId string, addonName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetesAddonStatus(clusterId, addonName)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return object, "Deleted", nil
			}
			return nil, "", WrapError(err)
		}
		if object.Version == "" {
			return object, "Deleted", nil
		}

		return object, "Running", nil
	}
}

func (s *CsClient) installAddon(d *schema.ResourceData) error {
	clusterId := d.Get("cluster_id").(string)

	body := make([]*client.InstallClusterAddonsRequestBody, 0)
	b := &client.InstallClusterAddonsRequestBody{
		Name:    tea.String(d.Get("name").(string)),
		Version: tea.String(d.Get("version").(string)),
	}

	if config, exist := d.GetOk("config"); exist {
		b.Config = tea.String(config.(string))
	}
	body = append(body, b)

	creationArgs := &client.InstallClusterAddonsRequest{
		Body: body,
	}

	_, err := s.client.InstallClusterAddons(&clusterId, creationArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "installAddon", err)
	}

	return nil
}

func (s *CsClient) upgradeAddon(d *schema.ResourceData) error {
	clusterId := d.Get("cluster_id").(string)

	body := make([]*client.UpgradeClusterAddonsRequestBody, 0)
	b := &client.UpgradeClusterAddonsRequestBody{
		ComponentName: tea.String(d.Get("name").(string)),
		NextVersion:   tea.String(d.Get("version").(string)),
	}

	if config, exist := d.GetOk("config"); exist {
		b.Config = tea.String(config.(string))
	}

	body = append(body, b)

	upgradeArgs := &client.UpgradeClusterAddonsRequest{
		Body: body,
	}

	_, err := s.client.UpgradeClusterAddons(&clusterId, upgradeArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
	}

	return nil
}

func (s *CsClient) uninstallAddon(d *schema.ResourceData) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	clusterId := parts[0]

	body := make([]*client.UnInstallClusterAddonsRequestAddons, 0)
	b := &client.UnInstallClusterAddonsRequestAddons{
		Name: tea.String(parts[1]),
	}
	body = append(body, b)

	uninstallArgs := &client.UnInstallClusterAddonsRequest{
		Addons: body,
	}

	_, err = s.client.UnInstallClusterAddons(&clusterId, uninstallArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "uninstallAddon", err)
	}

	return nil
}

func (s *CsClient) updateAddonConfig(d *schema.ResourceData) error {
	clusterId := d.Get("cluster_id").(string)
	ComponentName := d.Get("name").(string)

	upgradeArgs := &client.ModifyClusterAddonRequest{
		Config: tea.String(d.Get("config").(string)),
	}

	_, err := s.client.ModifyClusterAddon(&clusterId, &ComponentName, upgradeArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddonConfig", err)
	}

	return nil
}

// This function returns the status of all available addons of the cluster
func (s *CsClient) DescribeCsKubernetesAddonMetadata(clusterId string, name string, version string) (*Component, error) {
	resp, err := s.client.DescribeClusterAddonMetadata(&clusterId, &name, &version)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesExistedAddons", err)
	}
	result := &Component{
		ComponentName: *resp.Body.Name,
		Version:       *resp.Body.Version,
		ConfigSchema:  *resp.Body.ConfigSchema,
	}

	return result, nil
}

func (s *CsService) DescribeCsKubernetesNodePool(id string) (nodePool *cs.NodePoolDetail, err error) {
	invoker := NewInvoker()
	var requestInfo *cs.Client
	var response interface{}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	clusterId := parts[0]
	nodePoolId := parts[1]

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeNodePoolDetail(clusterId, nodePoolId)
		})
		response = raw
		return err
	}); err != nil {
		if e, ok := err.(*common.Error); ok {
			for _, code := range []int{400} {
				if e.StatusCode == code {
					return nil, WrapErrorf(err, NotFoundMsg, DenverdinoAliyungo)
				}
			}
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, nodePoolId, "DescribeNodePool", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = clusterId
		requestMap["NodePoolId"] = nodePoolId
		addDebug("DescribeNodepool", response, requestInfo, requestMap)
	}
	nodePool, _ = response.(*cs.NodePoolDetail)
	if nodePool.NodePoolId != nodePoolId {
		return nil, WrapErrorf(Error(GetNotFoundMessage("CsNodePool", nodePoolId)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *CsService) WaitForCsKubernetes(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCsKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ClusterId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
}

func (s *CsService) DescribeCsManagedKubernetes(id string) (cluster *cs.KubernetesClusterDetail, err error) {
	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeKubernetesClusterDetail(id)
		})
		response = raw
		return err
	}); err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, "DescribeKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = id
		addDebug("DescribeKubernetesCluster", response, requestInfo, requestMap, map[string]interface{}{"Id": id})
	}
	cluster, _ = response.(*cs.KubernetesClusterDetail)
	if cluster.ClusterId != id {
		return cluster, WrapErrorf(Error(GetNotFoundMessage("CSManagedKubernetes", id)), NotFoundMsg, ProviderERROR)
	}
	return

}

func (s *CsService) WaitForCSManagedKubernetes(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCsManagedKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ClusterId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
}

func (s *CsService) CsKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) CsKubernetesNodePoolStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetesNodePool(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) CsManagedKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsManagedKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) CsServerlessKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsServerlessKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) DescribeCsServerlessKubernetes(id string) (*cs.ServerlessClusterResponse, error) {
	cluster := &cs.ServerlessClusterResponse{}
	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeServerlessKubernetesCluster(id)
		})
		response = raw
		return err
	}); err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, WrapErrorf(err, NotFoundMsg, DenverdinoAliyungo)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, "DescribeServerlessKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = id
		addDebug("DescribeServerlessKubernetesCluster", response, requestInfo, requestMap, map[string]interface{}{"Id": id})
	}
	cluster, _ = response.(*cs.ServerlessClusterResponse)
	if cluster != nil && cluster.ClusterId != id {
		return cluster, WrapErrorf(Error(GetNotFoundMessage("CSServerlessKubernetes", id)), NotFoundMsg, ProviderERROR)
	}
	return cluster, nil

}

func (s *CsService) WaitForCSServerlessKubernetes(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCsServerlessKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ClusterId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
}

func (s *CsService) tagsToMap(tags []cs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *CsService) ignoreTag(t cs.Tag) bool {
	filter := []string{"^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *CsService) GetPermanentToken(clusterId string) (string, error) {

	describeClusterTokensResponse, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return csClient.DescribeClusterTokens(clusterId)
	})
	if err != nil {
		return "", WrapError(fmt.Errorf("failed to get permanent token,because of %v", err))
	}

	tokens, ok := describeClusterTokensResponse.([]*cs.ClusterTokenResponse)

	if ok != true {
		return "", WrapError(fmt.Errorf("failed to parse ClusterTokenResponse of cluster %s", clusterId))
	}

	permanentTokens := make([]string, 0)

	for _, token := range tokens {
		if token.Expired == 0 && token.IsActive == 1 {
			permanentTokens = append(permanentTokens, token.Token)
			break
		}
	}

	// create a new token
	if len(permanentTokens) == 0 {
		createClusterTokenResponse, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			clusterTokenReqeust := &cs.ClusterTokenReqeust{}
			clusterTokenReqeust.IsPermanently = true
			return csClient.CreateClusterToken(clusterId, clusterTokenReqeust)
		})
		if err != nil {
			return "", WrapError(fmt.Errorf("failed to create permanent token,because of %v", err))
		}

		token, ok := createClusterTokenResponse.(*cs.ClusterTokenResponse)
		if ok != true {
			return "", WrapError(fmt.Errorf("failed to parse token of %s", clusterId))
		}
		return token.Token, nil
	}

	return permanentTokens[0], nil
}

// GetUserData of cluster
func (s *CsService) GetUserData(clusterId string, labels string, taints string) (string, error) {

	token, err := s.GetPermanentToken(clusterId)

	if err != nil {
		return "", err
	}

	if labels == "" {
		labels = fmt.Sprintf("%s=true", DefaultECSTag)
	} else {
		labels = fmt.Sprintf("%s,%s=true", labels, DefaultECSTag)
	}

	cluster, err := s.DescribeCsKubernetes(clusterId)

	if err != nil {
		return "", WrapError(fmt.Errorf("failed to describe cs kuberentes cluster,because of %v", err))
	}

	extra_options := make([]string, 0)

	if len(labels) > 0 || len(taints) > 0 {

		if len(labels) != 0 {
			extra_options = append(extra_options, fmt.Sprintf("--labels %s", labels))
		}

		if len(taints) != 0 {
			extra_options = append(extra_options, fmt.Sprintf("--taints %s", taints))
		}
	}

	if network, err := GetKubernetesNetworkName(cluster); err == nil && network != "" {
		extra_options = append(extra_options, fmt.Sprintf("--network %s", network))
	}

	extra_options_in_line := strings.Join(extra_options, " ")

	version := cluster.CurrentVersion
	region := cluster.RegionId

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(ATTACH_SCRIPT_WITH_VERSION+extra_options_in_line, region, region, version, token))), nil
}

func (s *CsService) UpgradeCluster(clusterId string, args *cs.UpgradeClusterArgs) error {
	invoker := NewInvoker()
	err := invoker.Run(func() error {
		_, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.UpgradeCluster(clusterId, args)
		})
		if e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		return WrapError(err)
	}

	state, upgradeError := s.WaitForUpgradeCluster(clusterId, "Upgrade")
	if state == cs.Task_Status_Success && upgradeError == nil {
		return nil
	}

	// if upgrade failed cancel the task
	err = invoker.Run(func() error {
		_, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.CancelUpgradeCluster(clusterId)
		})
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return WrapError(upgradeError)
	}

	if state, err := s.WaitForUpgradeCluster(clusterId, "CancelUpgrade"); err != nil || state != cs.Task_Status_Success {
		log.Printf("[WARN] %s ACK Cluster cancel upgrade error: %#v", clusterId, err)
	}

	return WrapError(upgradeError)
}

func (s *CsService) WaitForUpgradeCluster(clusterId string, action string) (string, error) {
	err := resource.Retry(UpgradeClusterTimeout, func() *resource.RetryError {
		resp, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.QueryUpgradeClusterResult(clusterId)
		})
		if err != nil || resp == nil {
			return resource.RetryableError(err)
		}

		upgradeResult := resp.(*cs.UpgradeClusterResult)
		if upgradeResult.UpgradeStep == cs.UpgradeStep_Success {
			return nil
		}

		if upgradeResult.UpgradeStep == cs.UpgradeStep_Pause && upgradeResult.UpgradeStatus.Failed == "true" {
			msg := ""
			events := upgradeResult.UpgradeStatus.Events
			if len(events) > 0 {
				msg = events[len(events)-1].Message
			}
			return resource.NonRetryableError(fmt.Errorf("faild to %s cluster, error: %s", action, msg))
		}
		return resource.RetryableError(fmt.Errorf("%s cluster state not matched", action))
	})

	if err == nil {
		log.Printf("[INFO] %s ACK Cluster %s successed", action, clusterId)
		return cs.Task_Status_Success, nil
	}

	return cs.Task_Status_Failed, WrapError(err)
}

func GetKubernetesNetworkName(cluster *cs.KubernetesClusterDetail) (network string, err error) {

	metadata := make(map[string]interface{})
	if err := json.Unmarshal([]byte(cluster.MetaData), &metadata); err != nil {
		return "", fmt.Errorf("unmarshal metaData failed. error: %s", err)
	}

	for _, name := range NETWORK_ADDON_NAMES {
		if _, ok := metadata[fmt.Sprintf("%s%s", name, "Version")]; ok {
			return name, nil
		}
	}
	return "", fmt.Errorf("no network addon found")
}

func (s *CsClient) DescribeUserPermission(uid string) ([]*client.DescribeUserPermissionResponseBody, error) {
	body, err := s.client.DescribeUserPermission(tea.String(uid))
	if err != nil {
		return nil, err
	}

	return body.Body, err
}

func (s *CsClient) DescribeCsAutoscalingConfig(id string) (*client.CreateAutoscalingConfigRequest, error) {

	request := &client.CreateAutoscalingConfigRequest{
		CoolDownDuration:        tea.String("10m"),
		UnneededDuration:        tea.String("10m"),
		UtilizationThreshold:    tea.String("0.5"),
		GpuUtilizationThreshold: tea.String("0.5"),
		ScanInterval:            tea.String("30s"),
	}

	return request, nil
}

func (s *CsClient) DescribeTaskInfo(taskId string) string {
	if taskId == "" {
		return ""
	}
	resp, err := s.client.DescribeTaskInfo(tea.String(taskId))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("[TASK FAILED!!!]\nDetails: %++v", resp.Body.GoString())
}

func (s *CsClient) ModifyNodePoolNodeConfig(clusterId, nodepoolId string, request *client.ModifyNodePoolNodeConfigRequest) (interface{}, error) {
	log.Printf("[DEBUG] modifyNodePoolKubeletRequest %++v", *request)

	resp, err := s.client.ModifyNodePoolNodeConfig(tea.String(clusterId), tea.String(nodepoolId), request)
	if err != nil {
		return nil, WrapError(err)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = clusterId
		requestMap["NodePoolId"] = nodepoolId
		requestMap["Args"] = request
		addDebug("ModifyNodePoolKubeletConfig", resp, requestMap)
	}
	return resp, err
}

package alicloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	roacs "github.com/alibabacloud-go/cs-20151215/v3/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	aliyungoecs "github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const defaultNodePoolType = "ess"

func resourceAlicloudCSKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesNodePoolCreate,
		Read:   resourceAlicloudCSNodePoolRead,
		Update: resourceAlicloudCSNodePoolUpdate,
		Delete: resourceAlicloudCSNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instances", "desired_size"},
				Deprecated:    "Field 'node_count' has been deprecated from provider version 1.158.0. New field 'desired_size' instead.",
			},
			"desired_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instances", "node_count"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
			},
			"instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_name", "kms_encrypted_password"},
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "kms_encrypted_password"},
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"security_group_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'security_group_id' has been deprecated from provider version 1.145.0. New field 'security_group_ids' instead",
			},
			"system_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(20, 32768),
			},
			"system_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: csNodepoolDiskPerformanceLevelDiffSuppressFunc,
			},
			"system_disk_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_disk_kms_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"aes-256", "sm4-128"}, false),
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AliyunLinux", "Windows", "CentOS", "WindowsCore"}, false),
				Deprecated:   "Field 'platform' has been deprecated from provider version 1.145.0. New field 'image_type' instead",
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cpu_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "static"}, false),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12, 24, 36, 48, 60}),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Month"}, false),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"unschedulable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"data_disks": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"labels": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"taints": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"node_name_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"management": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_repair": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auto_upgrade": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"surge": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"surge_percentage": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"max_unavailable": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"scaling_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"max_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"cpu", "gpu", "gpushare", "spot"}, false),
						},
						"is_bond_eip": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_charge_type"},
						},
						"eip_internet_charge_type": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
							ConflictsWith: []string{"internet_charge_type"},
						},
						"eip_bandwidth": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 500),
							ConflictsWith: []string{"internet_charge_type"},
						},
					},
				},
				ConflictsWith: []string{"instances"},
			},
			"scaling_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"release", "recycle"}, false),
				DiffSuppressFunc: csNodepoolScalingPolicyDiffSuppressFunc,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByTraffic", "PayByBandwidth"}, false),
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"spot_strategy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"spot_price_limit": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"price_limit": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: csNodepoolSpotInstanceSettingDiffSuppressFunc,
			},
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:      100,
				ConflictsWith: []string{"node_count", "scaling_config", "desired_size"},
			},
			"keep_instance_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"format_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 5,
				Computed: true,
			},
			"image_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AliyunLinux", "AliyunLinux3", "AliyunLinux3Arm64", "AliyunLinuxUEFI", "CentOS", "Windows", "WindowsCore", "AliyunLinux Qboot", "ContainerOS"}, false),
			},
			"runtime_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"runtime_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cis_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"soc_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"polardb_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kubelet_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"registry_pull_qps": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"registry_burst": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"event_record_qps": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"event_burst": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_api_qps": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_api_burst": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"serialize_image_pulls": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu_manager_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eviction_hard": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"eviction_soft": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"eviction_soft_grace_period": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"system_reserved": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"kube_reserved": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
			"rollout_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_unavailable": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				Deprecated: "Field 'rollout_policy' has been deprecated from provider version 1.184.0. Please use new field 'rolling_policy' instead it to ensure the config takes effect",
			},
			"rolling_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_parallelism": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCSKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var requestInfo *cs.Client
	var raw interface{}

	clusterId := d.Get("cluster_id").(string)
	// prepare args and set default value
	args, err := buildNodePoolArgs(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "PrepareKubernetesNodePoolArgs", err)
	}

	if err = invoker.Run(func() error {
		raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.CreateNodePool(args, d.Get("cluster_id").(string))
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "CreateKubernetesNodePool", raw)
	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = args
		addDebug("CreateKubernetesNodePool", raw, requestInfo, requestMap)
	}

	nodePool, ok := raw.(*cs.CreateNodePoolResponse)
	if ok != true {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "ParseKubernetesNodePoolResponse", raw)
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, nodePool.NodePoolID))

	// reset interval to 10s
	stateConf := BuildStateConf([]string{"initial", "scaling"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, "ResourceID:%s , TaskID:%s ", d.Id(), nodePool.TaskID)
	}

	// attach existing node
	if v, ok := d.GetOk("instances"); ok && v != nil {
		attachExistingInstance(d, meta)
	}

	return resourceAlicloudCSNodePoolRead(d, meta)
}

func resourceAlicloudCSNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	vpcService := VpcService{client}
	d.Partial(true)
	update := false
	invoker := NewInvoker()
	args := &cs.UpdateNodePoolRequest{
		RegionId:         common.Region(client.RegionId),
		NodePoolInfo:     cs.NodePoolInfo{},
		ScalingGroup:     cs.ScalingGroup{},
		KubernetesConfig: cs.KubernetesConfig{},
		AutoScaling:      cs.AutoScaling{},
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("node_count") {
		oldV, newV := d.GetChange("node_count")
		oldValue, ok := oldV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count new value can not be parsed"), "parseError %d", newValue)
		}

		_, exist := d.GetOk("desired_size")
		if exist {
			update = true
			desiredSize := int64(newValue)
			args.ScalingGroup.DesiredSize = &desiredSize
		} else {
			if newValue < oldValue {
				removeNodePoolNodes(d, meta, parts, nil, nil)
				// The removal of a node is logically independent.
				// The removal of a node should not involve parameter changes.
				return resourceAlicloudCSNodePoolRead(d, meta)
			}
			update = true
			args.Count = int64(newValue) - int64(oldValue)
		}
	}

	if d.HasChange("name") {
		update = true
		args.NodePoolInfo.Name = tea.TransInterfaceToString(d.Get("name"))
	}
	if d.HasChange("vswitch_ids") {
		update = true
		var vswitchID string
		if list := expandStringList(d.Get("vswitch_ids").([]interface{})); len(list) > 0 {
			vswitchID = list[0]
		} else {
			vswitchID = ""
		}

		var vpcId string
		if vswitchID != "" {
			vsw, err := vpcService.DescribeVSwitch(vswitchID)
			if err != nil {
				return err
			}
			vpcId = vsw.VpcId
		}
		args.ScalingGroup.VpcId = vpcId
		args.ScalingGroup.VswitchIds = expandStringList(d.Get("vswitch_ids").([]interface{}))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		args.InstanceChargeType = tea.TransInterfaceToString(v)
		if tea.StringValue(args.InstanceChargeType) == string(PrePaid) {
			update = true
			args.Period = tea.TransInterfaceToInt(d.Get("period"))
			args.PeriodUnit = tea.TransInterfaceToString(d.Get("period_unit"))
			args.AutoRenew = tea.TransInterfaceToBool(d.Get("auto_renew"))
			args.AutoRenewPeriod = tea.TransInterfaceToInt(d.Get("auto_renew_period"))
		}
	}

	if d.HasChange("image_type") {
		update = true
		args.ScalingGroup.ImageType = tea.TransInterfaceToString(d.Get("image_type"))
	}

	if d.HasChange("platform") {
		update = true
		args.ScalingGroup.Platform = tea.TransInterfaceToString(d.Get("platform"))
	}

	if d.HasChange("desired_size") {
		update = true
		size := int64(d.Get("desired_size").(int))
		args.ScalingGroup.DesiredSize = &size
	}

	if d.HasChange("install_cloud_monitor") {
		update = true
		args.CmsEnabled = tea.TransInterfaceToBool(d.Get("install_cloud_monitor"))
	}

	if d.HasChange("instance_types") {
		update = true
		args.ScalingGroup.InstanceTypes = expandStringList(d.Get("instance_types").([]interface{}))
	}

	args.ScalingGroup.SystemDiskEncrypted = tea.TransInterfaceToBool(d.Get("system_disk_encrypted"))

	if tea.BoolValue(args.ScalingGroup.SystemDiskEncrypted) && d.HasChanges("system_disk_encrypt_algorithm") {
		update = true
		args.ScalingGroup.SystemDiskEncryptAlgorithm = tea.TransInterfaceToString(d.Get("system_disk_encrypt_algorithm"))
	}

	if tea.BoolValue(args.ScalingGroup.SystemDiskEncrypted) && d.HasChanges("system_disk_kms_key") {
		update = true
		args.ScalingGroup.SystemDiskKMSKeyId = tea.TransInterfaceToString(d.Get("system_disk_kms_key"))
	}

	if d.HasChanges("system_disk_snapshot_policy_id") {
		update = true
		args.ScalingGroup.WorkerSnapshotPolicyId = tea.TransInterfaceToString(d.Get("system_disk_snapshot_policy_id"))
	}

	if d.HasChange("password") {
		update = true
		args.ScalingGroup.LoginPassword = tea.TransInterfaceToString(d.Get("password"))
	}

	if d.HasChange("key_name") {
		update = true
		args.ScalingGroup.KeyPair = tea.TransInterfaceToString(d.Get("key_name"))
	}

	if d.HasChange("security_group_id") {
		update = true
		args.ScalingGroup.SecurityGroupId = d.Get("security_group_id").(string)
	}

	if d.HasChange("system_disk_category") {
		update = true
		args.ScalingGroup.SystemDiskCategory = aliyungoecs.DiskCategory(d.Get("system_disk_category").(string))
	}

	if d.HasChange("system_disk_size") {
		update = true
		args.ScalingGroup.SystemDiskSize = tea.Int64(int64(d.Get("system_disk_size").(int)))
	}

	if d.HasChange("system_disk_performance_level") {
		update = true
		args.SystemDiskPerformanceLevel = tea.TransInterfaceToString(d.Get("system_disk_performance_level"))
	}

	if d.HasChange("image_id") {
		update = true
		args.ScalingGroup.ImageId = tea.TransInterfaceToString(d.Get("image_id"))
	}

	if d.HasChange("data_disks") {
		update = true
		setNodePoolDataDisks(&args.ScalingGroup, d)
	}

	if d.HasChange("tags") {
		update = true
		setNodePoolTags(&args.ScalingGroup, d)
	}

	if d.HasChange("labels") {
		update = true
		setNodePoolLabels(&args.KubernetesConfig, d)
	}

	if d.HasChange("taints") {
		update = true
		setNodePoolTaints(&args.KubernetesConfig, d)
	}

	if d.HasChange("node_name_mode") {
		update = true
		args.KubernetesConfig.NodeNameMode = d.Get("node_name_mode").(string)
	}

	if d.HasChange("user_data") {
		update = true
		if v := d.Get("user_data").(string); v != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v)
			if base64DecodeError == nil {
				args.KubernetesConfig.UserData = tea.String(v)
			} else {
				args.KubernetesConfig.UserData = tea.String(base64.StdEncoding.EncodeToString([]byte(v)))
			}
		}
	}

	if d.HasChange("scaling_config") {
		update = true
		args.AutoScaling = setAutoScalingConfig(d.Get("scaling_config").([]interface{}))
	}

	if d.HasChange("management") {
		update = true
		args.Management = setManagedNodepoolConfig(d.Get("management").([]interface{}))
	}

	if d.HasChange("internet_max_bandwidth_out") {
		update = true
		args.InternetMaxBandwidthOut = tea.TransInterfaceToInt(d.Get("internet_max_bandwidth_out"))
	}

	if d.HasChange("scaling_policy") {
		update = true
		args.ScalingPolicy = tea.TransInterfaceToString(d.Get("scaling_policy"))
	}

	// spot
	if d.HasChange("spot_strategy") {
		update = true
		args.SpotStrategy = tea.TransInterfaceToString(d.Get("spot_strategy"))
	}
	if d.HasChange("spot_price_limit") {
		update = true
		args.SpotPriceLimit = setSpotPriceLimit(d.Get("spot_price_limit").([]interface{}))
	}
	if d.HasChange("rds_instances") {
		update = true
		args.RdsInstances = expandStringList(d.Get("rds_instances").([]interface{}))
	}
	if d.HasChange("polardb_ids") {
		update = true
		args.PolarDBIds = expandStringList(d.Get("polardb_ids").([]interface{}))
	}

	if update {
		var response interface{}
		if err := invoker.Run(func() error {
			var err error
			response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				resp, err := csClient.UpdateNodePool(parts[0], parts[1], args)
				return resp, err
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateKubernetesNodePool", DenverdinoAliyungo)
		}
		if debugOn() {
			resizeRequestMap := make(map[string]interface{})
			resizeRequestMap["ClusterId"] = parts[0]
			resizeRequestMap["NodePoolId"] = parts[1]
			resizeRequestMap["Args"] = args
			addDebug("UpdateKubernetesNodePool", response, resizeRequestMap)
		}

		stateConf := BuildStateConf([]string{"scaling", "updating", "removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("kubelet_configuration") {
		roaClient, err := client.NewRoaCsClient()
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "InitClient", AlibabaCloudSdkGoERROR)
		}
		csClient := CsClient{roaClient}
		kubeletConfig := &roacs.ModifyNodePoolNodeConfigRequestKubeletConfig{}
		rolling := &roacs.ModifyNodePoolNodeConfigRequestRollingPolicy{}

		if v, ok := d.GetOk("kubelet_configuration"); ok {
			if err = setKubeletConfigParamsForUpdate(kubeletConfig, v.([]interface{})); err != nil {
				return WrapError(err)
			}
		}

		if v, ok := d.GetOk("rolling_policy"); ok {
			if err = setRollingPolicy(rolling, v.([]interface{})); err != nil {
				return WrapError(err)
			}
		}

		modifyNodePoolKubeletRequest := &roacs.ModifyNodePoolNodeConfigRequest{
			KubeletConfig: kubeletConfig,
			RollingPolicy: rolling,
		}

		resp, err := csClient.ModifyNodePoolNodeConfig(parts[0], parts[1], modifyNodePoolKubeletRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "ModifyNodePoolKubeletConfig", AlibabaCloudSdkGoERROR, resp)
		}
		modifyNodePoolKubeletResp, _ := resp.(*roacs.ModifyNodePoolNodeConfigResponse)

		stateConf := BuildStateConf([]string{"scaling", "updating", "removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsgWithTask, d.Id(), csClient.DescribeTaskInfo(tea.StringValue(modifyNodePoolKubeletResp.Body.TaskId)))
		}
	}

	// attach or remove existing node
	if d.HasChange("instances") {
		rawOldValue, rawNewValue := d.GetChange("instances")
		oldValue, ok := rawOldValue.([]interface{})
		if ok != true {
			return WrapErrorf(fmt.Errorf("instances old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := rawNewValue.([]interface{})
		if ok != true {
			return WrapErrorf(fmt.Errorf("instances new value can not be parsed"), "parseError %d", oldValue)
		}

		if len(newValue) > len(oldValue) {
			attachExistingInstance(d, meta)
		} else {
			removeNodePoolNodes(d, meta, parts, oldValue, newValue)
		}
	}

	_ = resource.Retry(10*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Start retry fetch node pool info: %s", d.Id())
		nodePoolDetail, err := csService.DescribeCsKubernetesNodePool(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if nodePoolDetail.TotalNodes != d.Get("node_count").(int) && nodePoolDetail.TotalNodes != d.Get("desired_size").(int) {
			time.Sleep(20 * time.Second)
			return resource.RetryableError(Error("[ERROR] The number of nodes is inconsistent %s", d.Id()))
		}

		return resource.NonRetryableError(Error("[DEBUG] The number of nodes is the same"))
	})

	update = false
	d.Partial(false)
	return resourceAlicloudCSNodePoolRead(d, meta)
}

func resourceAlicloudCSNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	object, err := csService.DescribeCsKubernetesNodePool(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("node_count", object.TotalNodes)
	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_ids", object.VswitchIds)
	d.Set("instance_types", object.InstanceTypes)
	d.Set("key_name", object.KeyPair)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("system_disk_category", object.SystemDiskCategory)
	d.Set("system_disk_size", object.SystemDiskSize)
	d.Set("system_disk_performance_level", object.SystemDiskPerformanceLevel)
	d.Set("image_id", object.ImageId)
	d.Set("platform", object.Platform)
	d.Set("scaling_policy", object.ScalingPolicy)
	d.Set("node_name_mode", object.NodeNameMode)
	d.Set("user_data", object.UserData)
	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("unschedulable", object.Unschedulable)
	d.Set("instance_charge_type", object.InstanceChargeType)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("spot_strategy", object.SpotStrategy)
	d.Set("internet_charge_type", object.InternetChargeType)
	d.Set("internet_max_bandwidth_out", object.InternetMaxBandwidthOut)
	d.Set("install_cloud_monitor", object.CmsEnabled)
	d.Set("image_type", object.ScalingGroup.ImageType)
	d.Set("security_group_ids", object.ScalingGroup.SecurityGroupIds)
	d.Set("runtime_name", object.Runtime)
	d.Set("runtime_version", object.RuntimeVersion)
	d.Set("deployment_set_id", object.DeploymentSetId)
	d.Set("cis_enabled", object.CisEnabled)
	d.Set("soc_enabled", object.SocEnabled)

	if object.DesiredSize != nil {
		d.Set("desired_size", *object.DesiredSize)
	}

	if tea.StringValue(object.InstanceChargeType) == "PrePaid" {
		d.Set("period", object.Period)
		d.Set("period_unit", object.PeriodUnit)
		d.Set("auto_renew", object.AutoRenew)
		d.Set("auto_renew_period", object.AutoRenewPeriod)
	}

	d.Set("system_disk_encrypted", object.SystemDiskEncrypted)
	d.Set("system_disk_kms_key", object.SystemDiskKMSKeyId)
	d.Set("system_disk_encrypt_algorithm", object.SystemDiskEncryptAlgorithm)
	d.Set("system_disk_snapshot_policy_id", object.WorkerSnapshotPolicyId)

	d.Set("rds_instances", object.RdsInstances)
	d.Set("polardb_ids", object.PolarDBIds)

	if passwd, ok := d.GetOk("password"); ok && passwd.(string) != "" {
		d.Set("password", passwd)
	}

	if parts, err := ParseResourceId(d.Id(), 2); err != nil {
		return WrapError(err)
	} else {
		d.Set("cluster_id", string(parts[0]))
	}

	if err := d.Set("data_disks", flattenNodeDataDisksConfig(object.DataDisks)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("taints", flattenTaintsConfig(object.Taints)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("labels", flattenLabelsConfig(object.Labels)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return WrapError(err)
	}

	if tea.BoolValue(object.Management.Enable) {
		if err := d.Set("management", flattenManagementNodepoolConfig(&object.Management)); err != nil {
			return WrapError(err)
		}
	}

	if tea.BoolValue(object.AutoScaling.Enable) {
		if err := d.Set("scaling_config", flattenAutoScalingConfig(&object.AutoScaling)); err != nil {
			return WrapError(err)
		}
	}

	if err := d.Set("spot_price_limit", flattenSpotPriceLimit(object.SpotPriceLimit)); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudCSNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	// delete all nodes [deprecated]
	// removeNodePoolNodes(d, meta, parts, nil, nil)

	// force delete
	var response interface{}
	err = resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.ForceDeleteNodePool(parts[0], parts[1])
			})
			response = raw
			return err
		}); err != nil {
			return resource.RetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = parts[0]
			requestMap["NodePoolId"] = parts[1]
			addDebug("DeleteClusterNodePool", response, d.Id(), requestMap)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNodePoolNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteClusterNodePool", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"active", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildNodePoolArgs(d *schema.ResourceData, meta interface{}) (*cs.CreateNodePoolRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	vpcService := VpcService{client}

	var vswitchID string
	if list := expandStringList(d.Get("vswitch_ids").([]interface{})); len(list) > 0 {
		vswitchID = list[0]
	} else {
		vswitchID = ""
	}

	var vpcId string
	if vswitchID != "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp
		}
	}

	creationArgs := &cs.CreateNodePoolRequest{
		RegionId: common.Region(client.RegionId),
		NodePoolInfo: cs.NodePoolInfo{
			Name:         tea.TransInterfaceToString(d.Get("name")),
			NodePoolType: defaultNodePoolType, // hard code the type
		},
		ScalingGroup: cs.ScalingGroup{
			VpcId:              vpcId,
			VswitchIds:         expandStringList(d.Get("vswitch_ids").([]interface{})),
			InstanceTypes:      expandStringList(d.Get("instance_types").([]interface{})),
			LoginPassword:      tea.String(password),
			KeyPair:            tea.TransInterfaceToString(d.Get("key_name")),
			SystemDiskCategory: aliyungoecs.DiskCategory(d.Get("system_disk_category").(string)),
			SystemDiskSize:     tea.Int64(int64(d.Get("system_disk_size").(int))),
			SecurityGroupId:    d.Get("security_group_id").(string),
			ImageId:            tea.TransInterfaceToString(d.Get("image_id")),
		},
		KubernetesConfig: cs.KubernetesConfig{
			NodeNameMode: d.Get("node_name_mode").(string),
		},
	}

	if v, ok := d.GetOkExists("node_count"); ok {
		creationArgs.Count = int64(v.(int))
	}

	if v, ok := d.GetOkExists("desired_size"); ok {
		size := int64(v.(int))
		creationArgs.DesiredSize = &size
	}

	setNodePoolDataDisks(&creationArgs.ScalingGroup, d)
	setNodePoolTags(&creationArgs.ScalingGroup, d)
	setNodePoolTaints(&creationArgs.KubernetesConfig, d)
	setNodePoolLabels(&creationArgs.KubernetesConfig, d)

	if v, ok := d.GetOk("instance_charge_type"); ok {
		creationArgs.InstanceChargeType = tea.TransInterfaceToString(v)
		if tea.StringValue(creationArgs.InstanceChargeType) == string(PrePaid) {
			creationArgs.Period = tea.TransInterfaceToInt(d.Get("period"))
			creationArgs.PeriodUnit = tea.TransInterfaceToString(d.Get("period_unit"))
			creationArgs.AutoRenew = tea.TransInterfaceToBool(d.Get("auto_renew"))
			creationArgs.AutoRenewPeriod = tea.TransInterfaceToInt(d.Get("auto_renew_period"))
		}
	}

	if v, ok := d.GetOkExists("system_disk_encrypted"); ok {
		creationArgs.SystemDiskEncrypted = tea.TransInterfaceToBool(v)
		if tea.BoolValue(creationArgs.SystemDiskEncrypted) {
			creationArgs.SystemDiskKMSKeyId = tea.TransInterfaceToString(d.Get("system_disk_kms_key"))
			creationArgs.SystemDiskEncryptAlgorithm = tea.TransInterfaceToString(d.Get("system_disk_encrypt_algorithm"))
		}
	}

	if v, ok := d.GetOk("deployment_set_id"); ok {
		creationArgs.DeploymentSetId = tea.TransInterfaceToString(v)
	}

	if v, ok := d.GetOk("install_cloud_monitor"); ok {
		creationArgs.CmsEnabled = tea.TransInterfaceToBool(v)
	}

	if v, ok := d.GetOk("unschedulable"); ok {
		creationArgs.Unschedulable = tea.TransInterfaceToBool(v)
	}

	if v, ok := d.GetOk("user_data"); ok && v != "" {
		data := v.(string)
		_, base64DecodeError := base64.StdEncoding.DecodeString(data)
		if base64DecodeError == nil {
			creationArgs.KubernetesConfig.UserData = tea.String(data)
		} else {
			creationArgs.KubernetesConfig.UserData = tea.String(base64.StdEncoding.EncodeToString([]byte(data)))
		}
	}

	// set auto scaling config
	if v, ok := d.GetOk("scaling_policy"); ok {
		creationArgs.ScalingPolicy = tea.TransInterfaceToString(v)
	}

	if v, ok := d.GetOk("scaling_config"); ok {
		if sc, ok := v.([]interface{}); len(sc) > 0 && ok {
			creationArgs.AutoScaling = setAutoScalingConfig(sc)
		}
	}

	// set manage nodepool params
	if v, ok := d.GetOk("management"); ok {
		if management, ok := v.([]interface{}); len(management) > 0 && ok {
			creationArgs.Management = setManagedNodepoolConfig(management)
		}
	}

	if v, ok := d.GetOk("system_disk_performance_level"); ok {
		creationArgs.SystemDiskPerformanceLevel = tea.TransInterfaceToString(v)
	}

	if v, ok := d.GetOk("system_disk_snapshot_policy_id"); ok {
		creationArgs.WorkerSnapshotPolicyId = tea.TransInterfaceToString(v)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		creationArgs.ResourceGroupId = tea.TransInterfaceToString(v)
	}

	// setting spot instance
	if v, ok := d.GetOk("spot_strategy"); ok {
		creationArgs.SpotStrategy = tea.TransInterfaceToString(v)
	}

	if v, ok := d.GetOk("spot_price_limit"); ok {
		creationArgs.SpotPriceLimit = setSpotPriceLimit(v.([]interface{}))
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		creationArgs.InternetChargeType = tea.TransInterfaceToString(v)
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		creationArgs.InternetMaxBandwidthOut = tea.TransInterfaceToInt(v)
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		creationArgs.SecurityGroupIds = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("platform"); ok {
		creationArgs.Platform = tea.TransInterfaceToString(v)
	}

	if v, ok := d.GetOk("image_type"); ok {
		creationArgs.ImageType = tea.TransInterfaceToString(v)
	}

	if v, ok := d.GetOk("runtime_name"); ok {
		creationArgs.Runtime = v.(string)
	}

	if v, ok := d.GetOk("runtime_version"); ok {
		creationArgs.RuntimeVersion = v.(string)
	}

	cisEnabled, socEnabled := false, false
	if v, ok := d.GetOkExists("cis_enabled"); ok {
		cisEnabled = v.(bool)
	}
	if v, ok := d.GetOkExists("soc_enabled"); ok {
		socEnabled = v.(bool)
	}
	if cisEnabled && socEnabled {
		return creationArgs, fmt.Errorf("setting SOC and CIS together is not supported")
	} else if cisEnabled {
		creationArgs.CisEnabled = tea.Bool(cisEnabled)
	} else if socEnabled {
		creationArgs.SocEnabled = tea.Bool(socEnabled)
	}

	if v, ok := d.GetOk("rds_instances"); ok {
		creationArgs.RdsInstances = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("polardb_ids"); ok {
		creationArgs.PolarDBIds = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("cpu_policy"); ok {
		creationArgs.CpuPolicy = v.(string)
	}

	// kubelet
	if v, ok := d.GetOk("kubelet_configuration"); ok {
		config, err := setKubeletConfigParamsForCreate(v.([]interface{}))
		if err != nil {
			return creationArgs, WrapError(err)
		}
		creationArgs.NodeConfig = &cs.NodeConfig{}
		creationArgs.NodeConfig.KubeletConfiguration = config
	}

	return creationArgs, nil
}

func ConvertCsTags(d *schema.ResourceData) ([]cs.Tag, error) {
	tags := make([]cs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, cs.Tag{
						Key:   key,
						Value: v,
					})
				}
			}
		}
	}

	return tags, nil
}

func setNodePoolTags(scalingGroup *cs.ScalingGroup, d *schema.ResourceData) error {
	if tags, err := ConvertCsTags(d); err == nil {
		scalingGroup.Tags = tags
	}
	return nil
}

func setNodePoolLabels(config *cs.KubernetesConfig, d *schema.ResourceData) error {
	labels := make([]cs.Label, 0)
	if v, ok := d.GetOk("labels"); ok && len(v.([]interface{})) > 0 {
		vl := v.([]interface{})
		for _, i := range vl {
			if m, ok := i.(map[string]interface{}); ok {
				labels = append(labels, cs.Label{
					Key:   m["key"].(string),
					Value: m["value"].(string),
				})
			}
		}
	}
	config.Labels = labels

	return nil
}

func setNodePoolDataDisks(scalingGroup *cs.ScalingGroup, d *schema.ResourceData) error {
	if dds, ok := d.GetOk("data_disks"); ok {
		disks := dds.([]interface{})
		createDataDisks := make([]cs.NodePoolDataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := cs.NodePoolDataDisk{
				Size:                 pack["size"].(int),
				DiskName:             pack["name"].(string),
				Category:             pack["category"].(string),
				Device:               pack["device"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
				KMSKeyId:             pack["kms_key_id"].(string),
				Encrypted:            pack["encrypted"].(string),
				PerformanceLevel:     pack["performance_level"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		scalingGroup.DataDisks = createDataDisks
	}

	return nil
}

func setNodePoolTaints(config *cs.KubernetesConfig, d *schema.ResourceData) error {
	taints := make([]cs.Taint, 0)
	if v, ok := d.GetOk("taints"); ok && len(v.([]interface{})) > 0 {
		vl := v.([]interface{})
		for _, i := range vl {
			if m, ok := i.(map[string]interface{}); ok {
				taints = append(taints, cs.Taint{
					Key:    m["key"].(string),
					Value:  m["value"].(string),
					Effect: cs.Effect(m["effect"].(string)),
				})
			}

		}
	}
	config.Taints = taints

	return nil
}

func setManagedNodepoolConfig(l []interface{}) (config cs.Management) {
	if len(l) == 0 || l[0] == nil {
		config.Enable = tea.Bool(false)
		return config
	}

	m := l[0].(map[string]interface{})

	// Once "management" is set, we think of it as creating a managed node pool
	config.Enable = tea.Bool(true)

	if v, ok := m["auto_repair"].(bool); ok {
		config.AutoRepair = tea.Bool(v)
	}
	if v, ok := m["auto_upgrade"].(bool); ok {
		config.UpgradeConf.AutoUpgrade = tea.Bool(v)
	}
	if v, ok := m["surge"].(int); ok {
		config.UpgradeConf.Surge = tea.Int64(int64(v))
	}
	if v, ok := m["surge_percentage"].(int); ok {
		config.UpgradeConf.SurgePercentage = tea.Int64(int64(v))
	}
	if v, ok := m["max_unavailable"].(int); ok {
		config.UpgradeConf.MaxUnavailable = tea.Int64(int64(v))
	}

	return config
}

func setAutoScalingConfig(l []interface{}) (config cs.AutoScaling) {
	if len(l) == 0 || l[0] == nil {
		config.Enable = tea.Bool(false)
		return config
	}

	m := l[0].(map[string]interface{})

	// Once "scaling_config" is set, we think of it as creating a auto scaling node pool
	config.Enable = tea.Bool(true)

	if v, ok := m["min_size"].(int); ok {
		config.MinInstances = tea.Int64(int64(v))
	}
	if v, ok := m["max_size"].(int); ok {
		config.MaxInstances = tea.Int64(int64(v))
	}
	if v, ok := m["type"].(string); ok {
		config.Type = tea.String(v)
	}
	if v, ok := m["is_bond_eip"].(bool); ok {
		config.IsBindEip = tea.Bool(v)
	}
	if v, ok := m["eip_internet_charge_type"].(string); ok {
		config.EipInternetChargeType = tea.String(v)
	}
	if v, ok := m["eip_bandwidth"].(int); ok {
		config.EipBandWidth = tea.Int64(int64(v))
	}
	return config
}

func setSpotPriceLimit(l []interface{}) []cs.SpotPrice {
	config := make([]cs.SpotPrice, 0)
	if len(l) == 0 || l[0] == nil {
		return config
	}
	for _, v := range l {
		if m, ok := v.(map[string]interface{}); ok {
			config = append(config, cs.SpotPrice{
				InstanceType: m["instance_type"].(string),
				PriceLimit:   m["price_limit"].(string),
			})
		}
	}

	return config
}

func flattenSpotPriceLimit(config []cs.SpotPrice) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, spotInfo := range config {
		m = append(m, map[string]interface{}{
			"instance_type": spotInfo.InstanceType,
			"price_limit":   spotInfo.PriceLimit,
		})
	}

	return m
}

func flattenAutoScalingConfig(config *cs.AutoScaling) (m []map[string]interface{}) {
	if config == nil {
		return
	}
	m = append(m, map[string]interface{}{
		"min_size":                 config.MinInstances,
		"max_size":                 config.MaxInstances,
		"type":                     config.Type,
		"is_bond_eip":              config.IsBindEip,
		"eip_internet_charge_type": config.EipInternetChargeType,
		"eip_bandwidth":            config.EipBandWidth,
	})

	return
}

func flattenManagementNodepoolConfig(config *cs.Management) (m []map[string]interface{}) {
	if config == nil {
		return
	}
	m = append(m, map[string]interface{}{
		"auto_repair":      config.AutoRepair,
		"auto_upgrade":     config.UpgradeConf.AutoUpgrade,
		"surge":            config.UpgradeConf.Surge,
		"surge_percentage": config.UpgradeConf.SurgePercentage,
		"max_unavailable":  config.UpgradeConf.MaxUnavailable,
	})

	return
}

func flattenNodeDataDisksConfig(config []cs.NodePoolDataDisk) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, disks := range config {
		m = append(m, map[string]interface{}{
			"size":              disks.Size,
			"category":          disks.Category,
			"encrypted":         disks.Encrypted,
			"performance_level": disks.PerformanceLevel,
			"kms_key_id":        disks.KMSKeyId,
		})
	}

	return m
}

func flattenTaintsConfig(config []cs.Taint) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, taint := range config {
		m = append(m, map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		})
	}

	return m
}

func flattenLabelsConfig(config []cs.Label) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, label := range config {
		m = append(m, map[string]interface{}{
			"key":   label.Key,
			"value": label.Value,
		})
	}

	return m
}

func flattenTagsConfig(config []cs.Tag) map[string]string {
	m := make(map[string]string, len(config))
	if len(config) < 0 {
		return m
	}

	for _, tag := range config {
		if tag.Key != DefaultClusterTag {
			m[tag.Key] = tag.Value
		}
	}

	return m
}

func removeNodePoolNodes(d *schema.ResourceData, meta interface{}, parseId []string, oldNodes []interface{}, newNodes []interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var response interface{}
	// list all nodes of the nodepool
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			nodes, _, err := csClient.GetKubernetesClusterNodes(parseId[0], common.Pagination{PageNumber: 1, PageSize: PageSizeLarge}, parseId[1])
			return nodes, err
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
	}

	ret := response.([]cs.KubernetesNodeType)
	// fetch the NodeName of all nodes
	var allNodeName []string
	for _, value := range ret {
		allNodeName = append(allNodeName, value.NodeName)
	}

	removeNodesName := allNodeName

	// remove automatically created nodes
	if d.HasChange("node_count") {
		o, n := d.GetChange("node_count")
		count := o.(int) - n.(int)
		removeNodesName = allNodeName[:count]
	}

	// remove manually added nodes
	if d.HasChange("instances") {
		var removeInstanceList []string
		var attachNodeList []string
		if oldNodes != nil && newNodes != nil {
			attachNodeList = difference(expandStringList(oldNodes), expandStringList(newNodes))
		}
		if len(newNodes) == 0 {
			attachNodeList = expandStringList(oldNodes)
		}
		for _, v := range ret {
			for _, name := range attachNodeList {
				if name == v.InstanceId {
					removeInstanceList = append(removeInstanceList, v.NodeName)
				}
			}
		}
		removeNodesName = removeInstanceList
	}

	removeNodesArgs := &cs.DeleteKubernetesClusterNodesRequest{
		Nodes:       removeNodesName,
		ReleaseNode: true,
		DrainNode:   false,
	}
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			resp, err := csClient.DeleteKubernetesClusterNodes(parseId[0], removeNodesArgs)
			return resp, err
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteKubernetesClusterNodes", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.SetPartial("node_count")

	return nil
}

func attachExistingInstance(d *schema.ResourceData, meta interface{}) error {
	csService := CsService{meta.(*connectivity.AliyunClient)}
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	clusterId := parts[0]
	nodePoolId := parts[1]

	args := &roacs.AttachInstancesRequest{
		NodepoolId:       tea.String(nodePoolId),
		FormatDisk:       tea.Bool(false),
		KeepInstanceName: tea.Bool(true),
	}

	if v, ok := d.GetOk("password"); ok {
		args.Password = tea.String(v.(string))
	}

	if v, ok := d.GetOk("key_name"); ok {
		args.KeyPair = tea.String(v.(string))
	}

	if v, ok := d.GetOk("format_disk"); ok {
		args.FormatDisk = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("keep_instance_name"); ok {
		args.KeepInstanceName = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_id"); ok {
		args.ImageId = tea.String(v.(string))
	}

	if v, ok := d.GetOk("instances"); ok {
		args.Instances = tea.StringSlice(expandStringList(v.([]interface{})))
	}

	_, err = client.AttachInstances(tea.String(clusterId), args)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "AttachInstances", AliyunTablestoreGoSdk)
	}

	stateConf := BuildStateConf([]string{"scaling"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func setKubeletConfigParamsForUpdate(config *roacs.ModifyNodePoolNodeConfigRequestKubeletConfig, l []interface{}) error {
	if len(l) <= 0 || l[0] == nil {
		return nil
	}
	m := l[0].(map[string]interface{})

	var (
		intVal  int64
		boolVal bool
		err     error
	)

	if v, ok := m["registry_pull_qps"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return WrapError(fmt.Errorf("failed to parse 'registry_pull_qps' due to %v", err))
		}
		config.RegistryPullQPS = tea.Int64(intVal)
	}
	if v, ok := m["registry_burst"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return WrapError(fmt.Errorf("failed to parse 'registry_burst' due to %v", err))
		}
		config.RegistryBurst = tea.Int64(intVal)
	}
	if v, ok := m["event_record_qps"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return WrapError(fmt.Errorf("failed to parse 'event_record_qps' due to %v", err))
		}
		config.EventRecordQPS = tea.Int64(intVal)
	}
	if v, ok := m["event_burst"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return WrapError(fmt.Errorf("failed to parse 'event_burst' due to %v", err))
		}
		config.EventBurst = tea.Int64(intVal)
	}
	if v, ok := m["kube_api_qps"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return WrapError(fmt.Errorf("failed to parse 'kube_api_qps' due to %v", err))
		}
		config.KubeAPIQPS = tea.Int64(intVal)
	}
	if v, ok := m["kube_api_burst"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return WrapError(fmt.Errorf("failed to parse 'kube_api_burst' due to %v", err))
		}
		config.KubeAPIBurst = tea.Int64(intVal)
	}
	if v, ok := m["serialize_image_pulls"]; ok && reflect.ValueOf(v).String() != "" {
		if boolVal, err = strconv.ParseBool(v.(string)); err != nil {
			return WrapError(fmt.Errorf("failed to parse 'serialize_image_pulls' due to %v", err))
		}
		config.SerializeImagePulls = tea.Bool(boolVal)
	}
	if v, ok := m["cpu_manager_policy"]; ok && reflect.ValueOf(v).String() != "" {
		config.CpuManagerPolicy = tea.String(v.(string))
	}
	if v, ok := m["eviction_hard"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.EvictionHard = v.(map[string]interface{})
	}
	if v, ok := m["eviction_soft"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.EvictionSoft = v.(map[string]interface{})
	}
	if v, ok := m["eviction_soft_grace_period"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.EvictionSoftGracePeriod = v.(map[string]interface{})
	}
	if v, ok := m["system_reserved"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.SystemReserved = v.(map[string]interface{})
	}
	if v, ok := m["kube_reserved"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.KubeReserved = v.(map[string]interface{})
	}

	return nil
}

func setKubeletConfigParamsForCreate(l []interface{}) (*cs.KubeletConfiguration, error) {
	config := &cs.KubeletConfiguration{}
	if len(l) <= 0 || l[0] == nil {
		return nil, nil
	}

	var (
		intVal  int64
		boolVal bool
		err     error
	)

	m := l[0].(map[string]interface{})

	if v, ok := m["registry_pull_qps"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return config, WrapError(fmt.Errorf("failed to parse 'registry_pull_qps' due to %v", err))
		}
		config.RegistryPullQPS = tea.Int64(intVal)
	}
	if v, ok := m["registry_burst"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return config, WrapError(fmt.Errorf("failed to parse 'registry_burst' due to %v", err))
		}
		config.RegistryBurst = tea.Int64(intVal)
	}
	if v, ok := m["event_record_qps"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return config, WrapError(fmt.Errorf("failed to parse 'event_record_qps' due to %v", err))
		}
		config.EventRecordQPS = tea.Int64(intVal)
	}
	if v, ok := m["event_burst"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return config, WrapError(fmt.Errorf("failed to parse 'event_burst' due to %v", err))
		}
		config.EventBurst = tea.Int64(intVal)
	}
	if v, ok := m["kube_api_qps"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return config, WrapError(fmt.Errorf("failed to parse 'kube_api_qps' due to %v", err))
		}
		config.KubeAPIQPS = tea.Int64(intVal)
	}
	if v, ok := m["kube_api_burst"]; ok && reflect.ValueOf(v).String() != "" {
		if intVal, err = strconv.ParseInt(v.(string), 10, 64); err != nil {
			return config, WrapError(fmt.Errorf("failed to parse 'kube_api_burst' due to %v", err))
		}
		config.KubeAPIBurst = tea.Int64(intVal)
	}
	if v, ok := m["serialize_image_pulls"]; ok && reflect.ValueOf(v).String() != "" {
		if boolVal, err = strconv.ParseBool(v.(string)); err != nil {
			return config, WrapError(fmt.Errorf("failed to parse 'serialize_image_pulls' due to %v", err))
		}
		config.SerializeImagePulls = tea.Bool(boolVal)
	}
	if v, ok := m["cpu_manager_policy"]; ok && reflect.ValueOf(v).String() != "" {
		config.CpuManagerPolicy = tea.String(v.(string))
	}
	if v, ok := m["eviction_hard"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.EvictionHard = v.(map[string]interface{})
	}
	if v, ok := m["eviction_soft"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.EvictionSoft = v.(map[string]interface{})
	}
	if v, ok := m["eviction_soft_grace_period"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.EvictionSoftGracePeriod = v.(map[string]interface{})
	}
	if v, ok := m["system_reserved"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.SystemReserved = v.(map[string]interface{})
	}
	if v, ok := m["kube_reserved"]; ok && reflect.TypeOf(v).Kind() == reflect.Map {
		config.KubeReserved = v.(map[string]interface{})
	}
	return config, nil
}

func setRollingPolicy(policy *roacs.ModifyNodePoolNodeConfigRequestRollingPolicy, l []interface{}) error {
	if len(l) <= 0 || l[0] == nil {
		return nil
	}
	m := l[0].(map[string]interface{})
	if v, ok := m["max_parallelism"]; ok {
		policy.MaxParallelism = tea.Int64(int64(v.(int)))
	}
	return nil
}

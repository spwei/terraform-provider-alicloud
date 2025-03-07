package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type R_kvstoreService struct {
	client *connectivity.AliyunClient
}
type RKvstoreService struct {
	client *connectivity.AliyunClient
}

func (s *R_kvstoreService) DescribeInstanceSSL(id string) (object r_kvstore.DescribeInstanceSSLResponse, err error) {
	request := r_kvstore.CreateDescribeInstanceSSLRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeInstanceSSL(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeInstanceSSLResponse)
	return *response, nil
}

func (s *R_kvstoreService) DescribeSecurityIps(id string) (object r_kvstore.SecurityIpGroup, err error) {
	request := r_kvstore.CreateDescribeSecurityIpsRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeSecurityIps(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeSecurityIpsResponse)

	if len(response.SecurityIpGroups.SecurityIpGroup) < 1 {
		return object, nil
	} else {
		for _, v := range response.SecurityIpGroups.SecurityIpGroup {
			if v.SecurityIpGroupName == "ali_dms_group" || v.SecurityIpGroupName == "hdm_security_ips" {
				continue
			}
			return v, nil
		}
	}
	return response.SecurityIpGroups.SecurityIpGroup[0], nil
}

func (s *R_kvstoreService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]r_kvstore.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, r_kvstore.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	// 对系统 Tag 进行过滤
	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !ignoredTags(v, "") {
			removedTagKeys = append(removedTagKeys, v)
		}
	}
	if len(removedTagKeys) > 0 {
		request := r_kvstore.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removedTagKeys
		raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDbInstanceId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := r_kvstore.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDbInstanceId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *R_kvstoreService) DescribeInstanceAutoRenewalAttribute(id string) (object r_kvstore.Item, err error) {
	request := r_kvstore.CreateDescribeInstanceAutoRenewalAttributeRequest()
	request.RegionId = s.client.RegionId

	request.DBInstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeInstanceAutoRenewalAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeInstanceAutoRenewalAttributeResponse)

	if len(response.Items.Item) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return object, err
	}
	return response.Items.Item[0], nil
}

func (s *R_kvstoreService) DescribeSecurityGroupConfiguration(id string) (object r_kvstore.ItemsInDescribeSecurityGroupConfiguration, err error) {
	request := r_kvstore.CreateDescribeSecurityGroupConfigurationRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeSecurityGroupConfiguration(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeSecurityGroupConfigurationResponse)

	if len(response.Items.EcsSecurityGroupRelation) < 1 {
		return object, nil
	}
	return response.Items, nil
}

func (s *R_kvstoreService) DescribeKvstoreInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstanceAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	v, err := jsonpath.Get("$.Instances.DBInstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances", response)
	}
	if v == nil || len(v.([]interface{})) < 1 || fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
		return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *R_kvstoreService) KvstoreInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["InstanceStatus"] == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *R_kvstoreService) KvstoreInstanceAttributeRefreshFunc(id, attribute string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		return object, fmt.Sprint(object[attribute]), nil
	}
}

func (s *R_kvstoreService) DescribeKvstoreConnection(id string) (object []r_kvstore.InstanceNetInfo, err error) {
	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeDBInstanceNetInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreConnection", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)

	if len(response.NetInfoItems.InstanceNetInfo) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreConnection", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.NetInfoItems.InstanceNetInfo, nil
}

func (s *R_kvstoreService) DescribeKvstoreAccount(id string) (object r_kvstore.Account, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := r_kvstore.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.AccountName = parts[1]
	request.InstanceId = parts[0]

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeAccounts(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreAccount", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeAccountsResponse)

	if len(response.Accounts.Account) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreAccount", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Accounts.Account[0], nil
}

func (s *R_kvstoreService) KvstoreAccountStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreAccount(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.AccountStatus == failState {
				return object, object.AccountStatus, WrapError(Error(FailedToReachTargetStatus, object.AccountStatus))
			}
		}
		return object, object.AccountStatus, nil
	}
}

func (s *R_kvstoreService) DescribeBackupPolicy(id string) (object r_kvstore.DescribeBackupPolicyResponse, err error) {
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeBackupPolicy(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeBackupPolicyResponse)
	return *response, nil
}

func (s *RKvstoreService) DescribeKvstoreAuditLogConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAuditLogConfig"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Redis:AuditLogConfig", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) DescribeInstanceAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstanceAttribute"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Redis:AuditLogConfig", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Instances.DBInstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances.DBInstanceAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) KvstoreAuditLogConfigStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstanceAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["InstanceStatus"]) == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *RKvstoreService) DescribeInstanceAutoRenewalAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstanceAutoRenewalAttribute"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidOrderCharge.NotSupport"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.Item", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.Item", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) DescribeInstanceSSL(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstanceSSL"
	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"IncorrectEngineVersion"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *RKvstoreService) DescribeSecurityGroupConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeSecurityGroupConfiguration"
	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.EcsSecurityGroupRelation", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.EcsSecurityGroupRelation", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) DescribeSecurityIps(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeSecurityIps"
	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SecurityIpGroups.SecurityIpGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecurityIpGroups.SecurityIpGroup", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Redis", id)), NotFoundWithResponse, response)
	} else {
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["SecurityIpGroupName"]) != "default" {
				return v.(map[string]interface{}), nil
			}
		}
	}

	return object, nil
}

func (s *RKvstoreService) DescribeInstanceTDEStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstanceTDEStatus"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) DescribeEncryptionKey(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRedisaClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeEncryptionKey"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

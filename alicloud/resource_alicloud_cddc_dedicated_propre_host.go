// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCddcDedicatedPropreHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCddcDedicatedPropreHostCreate,
		Read:   resourceAliCloudCddcDedicatedPropreHostRead,
		Update: resourceAliCloudCddcDedicatedPropreHostUpdate,
		Delete: resourceAliCloudCddcDedicatedPropreHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ecs_class_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_count": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"disk_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"system_disk_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"sys_disk_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"data_disk_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"sys_disk_capacity": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"ecs_deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ecs_host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ecs_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ecs_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ecs_unique_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecs_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"alisql", "tair", "mssql", "mysql"}, false),
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"os_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_inherit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Monthly"}, false),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCddcDedicatedPropreHostCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateMyBase"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCddcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["DedicatedHostGroupId"] = d.Get("dedicated_host_group_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["Engine"] = d.Get("engine")
	request["SecurityGroupId"] = d.Get("security_group_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if v, ok := d.GetOk("os_password"); ok {
		request["OsPassword"] = v
	}
	if v, ok := d.GetOk("period_type"); ok {
		request["PeriodType"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}
	if v, ok := d.GetOk("password_inherit"); ok {
		request["PasswordInherit"] = v
	}
	if v, ok := d.GetOk("ecs_host_name"); ok {
		request["EcsHostName"] = v
	}
	if v, ok := d.GetOk("ecs_instance_name"); ok {
		request["EcsInstanceName"] = v
	}
	if v, ok := d.GetOk("ecs_deployment_set_id"); ok {
		request["EcsDeploymentSetId"] = v
	}
	if v, ok := d.GetOk("ecs_unique_suffix"); ok {
		request["EcsUniqueSuffix"] = v
	}
	request["ZoneId"] = d.Get("ecs_zone_id")
	if v, ok := d.GetOk("ecs_class_list"); ok {
		eCSClassListMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["instanceType"] = dataLoopTmp["instance_type"]
			dataLoopMap["sysDiskCapacity"] = dataLoopTmp["sys_disk_capacity"]
			dataLoopMap["nodeCount"] = 1
			dataLoopMap["sysDiskType"] = dataLoopTmp["sys_disk_type"]
			dataLoopMap["diskType"] = dataLoopTmp["disk_type"]
			dataLoopMap["diskCapacity"] = dataLoopTmp["disk_capacity"]
			dataLoopMap["diskCount"] = dataLoopTmp["disk_count"]
			dataLoopMap["dataDiskPerformanceLevel"] = dataLoopTmp["data_disk_performance_level"]
			dataLoopMap["systemDiskPerformanceLevel"] = dataLoopTmp["system_disk_performance_level"]
			eCSClassListMaps = append(eCSClassListMaps, dataLoopMap)
		}
		eCSClassListMapsJson, err := json.Marshal(eCSClassListMaps)
		if err != nil {
			return WrapError(err)
		}
		request["ECSClassList"] = string(eCSClassListMapsJson)
	}

	request["PayType"] = convertCddcPayTypeRequest(d.Get("payment_type").(string))
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cddc_dedicated_propre_host", action, AlibabaCloudSdkGoERROR)
	}

	groupName, _ := jsonpath.Get("$.OrderList.OrderList[0].DedicatedHostGroupName", response)
	ecsIdString, _ := jsonpath.Get("$.OrderList.OrderList[0].ECSInstanceIds", response)
	var ecsIds []string
	err = json.Unmarshal([]byte(ecsIdString.(string)), &ecsIds)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, fmt.Errorf("unmarshaling ECSInstanceIds got an error: %#v", err))
	}
	if len(ecsIds) != 1 {
		return WrapErrorf(err, FailedGetAttributeMsg, "$.OrderList.OrderList[0].ECSInstanceIds", response)
	}
	d.SetId(fmt.Sprintf("%v:%v", groupName, ecsIds[0]))

	cddcServiceV2 := CddcServiceV2{client}
	sysDiskCapacity := d.Get("ecs_class_list.0.sys_disk_capacity")
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(sysDiskCapacity)}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cddcServiceV2.CddcDedicatedPropreHostStateRefreshFunc(d.Id(), "$.EcsClasses[0].SysDiskCapacity", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCddcDedicatedPropreHostRead(d, meta)
}

func resourceAliCloudCddcDedicatedPropreHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcServiceV2 := CddcServiceV2{client}

	objectRaw, err := cddcServiceV2.DescribeCddcDedicatedPropreHost(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cddc_dedicated_propre_host DescribeCddcDedicatedPropreHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("engine", objectRaw["Engine"])
	d.Set("dedicated_host_group_id", objectRaw["DedicatedHostGroupId"])
	dedicatedInstances1RawObj, _ := jsonpath.Get("$.DedicatedInstances[*]", objectRaw)
	dedicatedInstances1Raw := make([]interface{}, 0)
	if dedicatedInstances1RawObj != nil {
		dedicatedInstances1Raw = dedicatedInstances1RawObj.([]interface{})
	}

	dedicatedInstancesChild1Raw := dedicatedInstances1Raw[0].(map[string]interface{})
	d.Set("ecs_deployment_set_id", dedicatedInstancesChild1Raw["DeploymentSetId"])
	d.Set("ecs_host_name", dedicatedInstancesChild1Raw["HostName"])
	d.Set("ecs_instance_name", dedicatedInstancesChild1Raw["InstanceName"])
	d.Set("ecs_zone_id", dedicatedInstancesChild1Raw["ZoneId"])
	d.Set("image_id", dedicatedInstancesChild1Raw["ImageId"])
	d.Set("key_pair_name", dedicatedInstancesChild1Raw["KeyPairName"])
	d.Set("payment_type", convertCddcDedicatedInstancesInstanceChargeTypeResponse(dedicatedInstancesChild1Raw["InstanceChargeType"]))
	d.Set("security_group_id", dedicatedInstancesChild1Raw["SecurityGroupIds"])
	d.Set("vswitch_id", dedicatedInstancesChild1Raw["VSwitchId"])
	d.Set("vpc_id", dedicatedInstancesChild1Raw["VpcId"])
	ecsClasses1Raw := objectRaw["EcsClasses"]
	ecsClassListMaps := make([]map[string]interface{}, 0)
	if ecsClasses1Raw != nil {
		for _, ecsClassesChild1Raw := range ecsClasses1Raw.([]interface{}) {
			ecsClassListMap := make(map[string]interface{})
			ecsClassesChild1Raw := ecsClassesChild1Raw.(map[string]interface{})
			ecsClassListMap["data_disk_performance_level"] = ecsClassesChild1Raw["DataDiskPerformanceLevel"]
			ecsClassListMap["disk_capacity"] = ecsClassesChild1Raw["DataDiskCapacity"]
			ecsClassListMap["disk_count"] = ecsClassesChild1Raw["DataDiskCount"]
			ecsClassListMap["disk_type"] = ecsClassesChild1Raw["DataDiskType"]
			ecsClassListMap["instance_type"] = ecsClassesChild1Raw["InstanceType"]
			ecsClassListMap["sys_disk_capacity"] = ecsClassesChild1Raw["SysDiskCapacity"]
			ecsClassListMap["sys_disk_type"] = ecsClassesChild1Raw["SysDiskType"]
			ecsClassListMap["system_disk_performance_level"] = ecsClassesChild1Raw["SystemDiskPerformanceLevel"]
			ecsClassListMaps = append(ecsClassListMaps, ecsClassListMap)
		}
	}
	d.Set("ecs_class_list", ecsClassListMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("dedicated_host_group_name", parts[0])
	d.Set("ecs_instance_id", parts[1])

	return nil
}

func resourceAliCloudCddcDedicatedPropreHostUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Dedicated Propre Host.")
	return nil
}

func resourceAliCloudCddcDedicatedPropreHostDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Dedicated Propre Host. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertCddcDedicatedInstancesInstanceChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PrePaid":
		return "Subscription"
	}
	return source
}
func convertCddcPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "PrePaid"
	}
	return source
}

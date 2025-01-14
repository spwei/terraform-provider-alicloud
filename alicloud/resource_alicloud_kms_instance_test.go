package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kms Instance. >>> Resource test cases, automatically generated.
// Case 4048
func TestAccAliCloudKmsInstance_basic4048(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence4048)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KmsInstanceSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "1",
					"key_num":         "1000",
					"secret_num":      "0",
					"spec":            "1000",
					"product_version": "3",
					"vpc_id":          "${local.vpc_id}",
					"zone_ids": []string{
						"cn-hangzhou-k", "cn-hangzhou-j"},
					"vswitch_ids": []string{
						"${local.vsw_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "1",
						"key_num":         "1000",
						"secret_num":      "0",
						"spec":            "1000",
						"product_version": "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vswitch.shareVswitch.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.shareVswitch.id}",
							"vpc_owner_id": "1511928242963727",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vswitch2.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vswitch2.id}",
							"vpc_owner_id": "1511928242963727",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vsw3.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vsw3.id}",
							"vpc_owner_id": "1511928242963727",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_num": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_num": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "vpc-bp14c07ucxg6h1xjmgcld",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "vsw-bp1wujtnspi1l3gvunvds",
							"vpc_owner_id": "1192853035118460",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "5",
					"key_num":         "2000",
					"secret_num":      "2000",
					"spec":            "2000",
					"renew_status":    "ManualRenewal",
					"product_version": "3",
					"renew_period":    "3",
					"vpc_id":          "${local.vpc_id}",
					"zone_ids": []string{
						"cn-hangzhou-k", "cn-hangzhou-j"},
					"vswitch_ids": []string{
						"${local.vsw_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "5",
						"key_num":         "2000",
						"secret_num":      "2000",
						"spec":            "2000",
						"renew_status":    "ManualRenewal",
						"product_version": "3",
						"renew_period":    "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status"},
			},
		},
	})
}

var AlicloudKmsInstanceMap4048 = map[string]string{
	"status":                   CHECKSET,
	"create_time":              CHECKSET,
	"instance_name":            CHECKSET,
	"ca_certificate_chain_pem": CHECKSET,
}

func AlicloudKmsInstanceBasicDependence4048(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "tf-testacc-kms-instance"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "default" {
  count = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
  cidr_block = "172.16.0.0/12"
  vpc_name   = "tf-testacc-kms-instance"
}

data "alicloud_vswitches" "vswitch" {
  vpc_id  = local.vpc_id
  zone_id = "cn-hangzhou-k"
}

data "alicloud_vswitches" "vswitch-j" {
  vpc_id  = local.vpc_id
  zone_id = "cn-hangzhou-j"
}

locals {
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids.0 : concat(alicloud_vpc.default[*].id, [""])[0]
  vsw_id = length(data.alicloud_vswitches.vswitch.ids) > 0 ? data.alicloud_vswitches.vswitch.ids.0 : concat(alicloud_vswitch.vswitch[*].id, [""])[0]
  vswj_id = length(data.alicloud_vswitches.vswitch-j.ids) > 0 ? data.alicloud_vswitches.vswitch-j.ids.0 : concat(alicloud_vswitch.vswitch-j[*].id, [""])[0]
}

resource "alicloud_vswitch" "vswitch" {
  count = length(data.alicloud_vswitches.vswitch.ids) > 0 ? 0 : 1
  vpc_id     = local.vpc_id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  count = length(data.alicloud_vswitches.vswitch-j.ids) > 0 ? 0 : 1
  vpc_id     = local.vpc_id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}3"
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}5"
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}7"
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}


`, name)
}

func AlicloudKmsInstanceBasicDependence4048_intl(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "tf-testacc-kms-instance"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "default" {
  count = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
  cidr_block = "172.16.0.0/12"
  vpc_name   = "tf-testacc-kms-instance"
}

data "alicloud_vswitches" "vswitch" {
  vpc_id  = local.vpc_id
  zone_id = "ap-southeast-1a"
}

data "alicloud_vswitches" "vswitch-j" {
  vpc_id  = local.vpc_id
  zone_id = "ap-southeast-1b"
}

locals {
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids.0 : concat(alicloud_vpc.default[*].id, [""])[0]
  vsw_id = length(data.alicloud_vswitches.vswitch.ids) > 0 ? data.alicloud_vswitches.vswitch.ids.0 : concat(alicloud_vswitch.vswitch[*].id, [""])[0]
  vswj_id = length(data.alicloud_vswitches.vswitch-j.ids) > 0 ? data.alicloud_vswitches.vswitch-j.ids.0 : concat(alicloud_vswitch.vswitch-j[*].id, [""])[0]
}

resource "alicloud_vswitch" "vswitch" {
  count = length(data.alicloud_vswitches.vswitch.ids) > 0 ? 0 : 1
  vpc_id     = local.vpc_id
  zone_id    = "ap-southeast-1a"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  count = length(data.alicloud_vswitches.vswitch-j.ids) > 0 ? 0 : 1
  vpc_id     = local.vpc_id
  zone_id    = "ap-southeast-1b"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}3"
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}5"
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}7"
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}


`, name)
}

// Case 4048  twin
func TestAccAliCloudKmsInstance_basic4048_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence4048)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KmsInstanceSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "7",
					"key_num":         "2000",
					"secret_num":      "1000",
					"spec":            "2000",
					"renew_status":    "ManualRenewal",
					"product_version": "3",
					"renew_period":    "3",
					"vpc_id":          "${local.vpc_id}",
					"zone_ids": []string{
						"cn-hangzhou-k", "cn-hangzhou-j"},
					"vswitch_ids": []string{
						"${local.vsw_id}"},
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vpc.shareVPC.id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.shareVswitch.id}",
							"vpc_owner_id": "1511928242963727",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vsw3.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vsw3.id}",
							"vpc_owner_id": "1511928242963727",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "7",
						"key_num":         "2000",
						"secret_num":      "1000",
						"spec":            "2000",
						"renew_status":    "ManualRenewal",
						"product_version": "3",
						"renew_period":    "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
						"bind_vpcs.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status"},
			},
		},
	})
}

func TestAccAlicloudKmsInstance_basic4048_intl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence4048_intl)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KmsInstanceIntlSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "2",
					"key_num":         "1000",
					"secret_num":      "1000",
					"spec":            "1000",
					"renew_status":    "ManualRenewal",
					"product_version": "3",
					"renew_period":    "3",
					"vpc_id":          "${local.vpc_id}",
					"zone_ids": []string{
						"ap-southeast-1a", "ap-southeast-1b"},
					"vswitch_ids": []string{
						"${local.vsw_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "2",
						"key_num":         "1000",
						"secret_num":      "1000",
						"spec":            "1000",
						"renew_status":    "ManualRenewal",
						"product_version": "3",
						"renew_period":    "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_num": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_num": "2000",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status"},
			},
		},
	})
}

// Test Kms Instance. <<< Resource test cases, automatically generated.

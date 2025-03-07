---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_acl_attachment"
sidebar_current: "docs-alicloud-resource-ga-acl-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Acl Attachment resource.
---

# alicloud_ga_acl_attachment

Provides a Global Accelerator (GA) Acl Attachment resource.

For information about Global Accelerator (GA) Acl Attachment and how to use it, see [What is Acl Attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-associateaclswithlistener).

-> **NOTE:** Available since v1.150.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  port_ranges {
    from_port = 80
    to_port   = 80
  }
}

resource "alicloud_ga_acl" "default" {
  acl_name           = "terraform-example"
  address_ip_version = "IPv4"
}

resource "alicloud_ga_acl_entry_attachment" "default" {
  acl_id            = alicloud_ga_acl.default.id
  entry             = "192.168.1.1/32"
  entry_description = "terraform-example"
}

resource "alicloud_ga_acl_attachment" "default" {
  acl_id      = alicloud_ga_acl.default.id
  listener_id = alicloud_ga_listener.default.id
  acl_type    = "white"
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of an ACL.
* `acl_type` - (Required, ForceNew) The type of the ACL. Valid values: `white`, `black`. 
  - `white`: Only requests from IP addresses or address segments in the selected access control list are forwarded. The whitelist applies to scenarios where applications only allow specific IP addresses. There are certain business risks in setting up a whitelist. Once the whitelist is set, only the IP addresses in the whitelist can access global acceleration listeners. If whitelist access is enabled, but no IP is added to the access policy group, the global acceleration listener will not forward the request.
  - `black`: All requests from IP addresses or address segments in the selected access control list are not forwarded. Blacklists are applicable to scenarios where applications restrict access to specific IP addresses. If blacklist access is enabled and no IP is added to the access policy group, the global acceleration listener forwards all requests.
* `dry_run` - (Optional) The dry run.
* `listener_id` - (Required, ForceNew) The ID of the listener.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Acl Attachment. The value formats as `<listener_id>:<acl_id>`.
* `status` - The status of the resource. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Acl Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Acl Attachment.

## Import

Global Accelerator (GA) Acl Attachment can be imported using the id. Format to `<listener_id>:<acl_id>`, e.g.

```shell
$ terraform import alicloud_ga_acl_attachment.example your_listener_id:your_acl_id
```

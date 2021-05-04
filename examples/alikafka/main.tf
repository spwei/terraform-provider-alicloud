data "alicloud_vpcs" "default" {
  is_default = true
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_alikafka_instance" "default" {
  name        = var.name
  topic_quota = var.topic_quota
  disk_type   = var.disk_type
  disk_size   = var.disk_size
  deploy_type = var.deploy_type
  io_max      = var.io_max
  eip_max     = var.eip_max
  paid_type   = var.paid_type
  spec_type   = var.spec_type
  vswitch_id  = var.vswitch_id != "" ? var.vswitch_id : data.alicloud_vswitches.default.ids.0
}
resource "alicloud_alikafka_consumer_group" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  consumer_id = var.consumer_id
}

resource "alicloud_alikafka_sasl_acl" "default" {
  instance_id               = alicloud_alikafka_instance.default.id
  username                  = alicloud_alikafka_sasl_user.default.username
  acl_resource_type         = var.acl_resource_type
  acl_resource_name         = var.acl_resource_type == "Topic" ? alicloud_alikafka_topic.default.topic : alicloud_alikafka_consumer_group.default.consumer_id
  acl_resource_pattern_type = var.acl_resource_pattern_type
  acl_operation_type        = var.acl_operation_type
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  username    = var.username
  password    = var.password
}

resource "alicloud_alikafka_topic" "default" {
  instance_id   = alicloud_alikafka_instance.default.id
  topic         = var.topic
  local_topic   = var.local_topic
  compact_topic = var.compact_topic
  partition_num = var.partition_num
  remark        = var.remark
}
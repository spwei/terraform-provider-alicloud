resource "alicloud_ons_instance" "instance" {
  instance_name = var.name
  remark        = "terraform-test-instance-remark"
}

resource "alicloud_ons_group" "default" {
  instance_id = alicloud_ons_instance.instance.id
  group_name  = var.group_name
  remark      = "terraform-test-group-remark"
}

resource "alicloud_ons_topic" "default" {
  instance_id  = alicloud_ons_instance.instance.id
  topic_name   = var.topic
  message_type = var.message_type
  remark       = var.topic_remark
}
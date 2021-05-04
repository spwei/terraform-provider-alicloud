variable "name" {
  default = "tf-testAccSagAclConfigName"
}
variable "sag_instance_id" {
  default = ""
}

resource "alicloud_sag_acl" "default" {
  name = var.name
}

resource "alicloud_sag_acl_rule" "default" {
  acl_id            = alicloud_sag_acl.default.id
  description       = "tf-testSagAclRule"
  policy            = "accept"
  ip_protocol       = "ALL"
  direction         = "in"
  source_cidr       = "10.10.10.0/24"
  source_port_range = "-1/-1"
  dest_cidr         = "192.168.10.0/24"
  dest_port_range   = "-1/-1"
  priority          = "1"
}

resource "alicloud_sag_client_user" "default" {
  count     = var.sag_instance_id == "" ? 0 : 1
  sag_id    = var.sag_instance_id
  bandwidth = "20"
  user_mail = "tftest-xxxxx@test.com"
  user_name = "th-username-xxxxx"
  password  = "xxxxxxx"
  client_ip = "192.1.10.0"
}
resource "alicloud_sag_dnat_entry" "default" {
  count         = var.sag_instance_id == "" ? 0 : 1
  sag_id        = var.sag_instance_id
  type          = "Intranet"
  ip_protocol   = "tcp"
  external_ip   = "1.0.0.2"
  external_port = "1"
  internal_ip   = "10.0.0.2"
  internal_port = "20"
}
resource "alicloud_sag_qos" "default" {
  name = var.name
}

resource "alicloud_sag_qos_car" "default" {
  qos_id            = alicloud_sag_qos.default.id
  name              = "tf-testSagQosCar"
  description       = "tf-testSagQosCar"
  priority          = "1"
  limit_type        = "Absolute"
  min_bandwidth_abs = "10"
  max_bandwidth_abs = "20"
}

resource "alicloud_sag_qos_policy" "default" {
  qos_id            = alicloud_sag_qos.default.id
  name              = "tf-testSagQosPolicy"
  description       = "tf-testSagQosPolicy"
  priority          = "1"
  ip_protocol       = "ALL"
  source_cidr       = "10.10.10.0/24"
  source_port_range = "-1/-1"
  dest_cidr         = "192.168.10.0/24"
  dest_port_range   = "-1/-1"
  start_time        = "2019-10-27T16:41:33+0800"
  end_time          = "2019-10-28T16:41:33+0800"
}

resource "alicloud_sag_snat_entry" "default" {
  count      = var.sag_instance_id == "" ? 0 : 1
  sag_id     = var.sag_instance_id
  cidr_block = "192.168.7.0/24"
  snat_ip    = "192.0.0.2"
}
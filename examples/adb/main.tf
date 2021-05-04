// Zones data source for availability_zone
data "alicloud_adb_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_adb_zones.default.zones.0.id
  vswitch_name = var.name
}
resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = var.db_cluster_category
  vswitch_id          = alicloud_vswitch.default.id
  mode                = var.mode
  compute_resource    = var.db_compute_resource
  payment_type        = "PayAsYouGo"
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF-update"
    For     = "acceptance-test-update"
  }
  security_ips = ["10.168.1.12", "10.168.1.11"]
}

resource "alicloud_nas_access_group" "main" {
  name        = "tf-testAccNasConfigName"
  type        = "Classic"
  description = "tf-testAccNasConfigDescription"
}

resource "alicloud_nas_access_rule" "main" {
  access_group_name = alicloud_nas_access_group.main.id
  source_cidr_ip    = "168.1.1.0/16"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 2
}

resource "alicloud_nas_file_system" "main" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "tf-testAccNasConfig"
  encrypt_type  = "1"
}

resource "alicloud_nas_mount_target" "main" {
  file_system_id    = alicloud_nas_file_system.main.id
  access_group_name = alicloud_nas_access_group.main.id
}


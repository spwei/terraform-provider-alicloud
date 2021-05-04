output "main_nas_access_group" {
  value = alicloud_nas_access_group.main.id
}

output "main-nas-accessrule" {
  value = alicloud_nas_access_rule.main.id
}

output "main_nas_file_system" {
  value = alicloud_nas_file_system.main.id
}

output "main-nas-mounttarget" {
  value = alicloud_nas_mount_target.main.id
}

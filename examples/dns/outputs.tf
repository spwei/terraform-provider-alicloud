output "domain" {
  value = alicloud_alidns_domain.dns.*.id
}

output "group" {
  value = alicloud_alidns_domain_group.group.*.id
}

output "record" {
  value = alicloud_alidns_record.record.*.id
}


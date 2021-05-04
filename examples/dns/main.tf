data "alicloud_alidns_domains" "domain" {
  domain_name_regex = "^*."
}

data "alicloud_alidns_records" "record" {
  domain_name = data.alicloud_alidns_domains.domain.domains[0].domain_name
  type        = "A"
  rr_regex    = "^@"
}

resource "alicloud_alidns_domain_group" "group" {
  domain_group_name = var.group_name
}

resource "alicloud_alidns_domain" "dns" {
  domain_name = var.domain_name
  group_id    = alicloud_alidns_domain_group.group.id
}

resource "alicloud_alidns_record" "record" {
  domain_name = alicloud_alidns_domain.dns.domain_name
  rr          = "alimail"
  type        = "CNAME"
  ttl         = 600
  priority    = 0
  value       = "mail.mxhichin.com"
  line        = "default"
  status      = "ENABLE"
  remark      = "test new domain record"
}



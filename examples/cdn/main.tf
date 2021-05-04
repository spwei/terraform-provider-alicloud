#####################
# Cdn
#####################
resource "alicloud_cdn_domain_new" "domain" {
  domain_name = var.domain_name
  cdn_type    = var.cdn_type
  scope       = var.scope
  dynamic "sources" {
    for_each = var.sources
    content {
      content  = lookup(sources.value, "content", null)
      type     = lookup(sources.value, "type", "ipaddr")
      port     = lookup(sources.value, "port", 80)
      priority = lookup(sources.value, "priority", 20)
      weight   = lookup(sources.value, "weight", 10)
    }
  }
  certificate_config {
    server_certificate = lookup(var.certificate_config[0], "server_certificate")
    private_key        = lookup(var.certificate_config[0], "private_key")
  }
  tags = merge({
    Name = var.domain_name
    }, var.tags
  )
}

#####################
# Cdn_config
#####################

resource "alicloud_cdn_domain_config" "this" {
  domain_name   = alicloud_cdn_domain_new.domain.domain_name
  function_name = var.function_name

  dynamic "function_args" {
    for_each = var.function_arg
    content {
      arg_name  = lookup(function_args.value, "arg_name", null)
      arg_value = lookup(function_args.value, "arg_value", null)
    }
  }
}
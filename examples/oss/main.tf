provider "alicloud" {
  alias  = "bj-prod"
  region = "cn-beijing"
}

resource "alicloud_oss_bucket" "bucket-new" {
  provider = alicloud.bj-prod

  bucket = var.bucket-new
  acl    = var.acl-bj
}

resource "alicloud_oss_bucket" "bucket-attr" {
  provider = alicloud.bj-prod

  bucket = var.bucket-attr

  website {
    index_document = var.index-doc
    error_document = var.error-doc
  }

  logging {
    target_bucket = alicloud_oss_bucket.bucket-new.id
    target_prefix = var.target-prefix
  }

  lifecycle_rule {
    id      = var.rule-days
    prefix  = "${var.rule-prefix}/${var.role-days}"
    enabled = true

    expiration {
      days = var.rule-days
    }
  }

  lifecycle_rule {
    id      = var.role-date
    prefix  = "${var.rule-prefix}/${var.role-date}"
    enabled = true

    expiration {
      date = var.rule-date
    }
  }

  referer_config {
    allow_empty = var.allow-empty
    referers    = [var.referers]
  }

  cors_rule {
    allowed_origins = [var.allow-origins-star]
    allowed_methods = split(",", var.allow-methods-put)
    allowed_headers = [var.allowed_headers]
  }

  cors_rule {
    allowed_origins = [var.allow-origins-aliyun]
    allowed_methods = split(",", var.allow-methods-get)
    allowed_headers = [var.allowed_headers]
    expose_headers  = [var.expose_headers]
    max_age_seconds = var.max_age_seconds
  }
}

resource "alicloud_oss_bucket_object" "content" {
  bucket  = alicloud_oss_bucket.bucket-new.bucket
  key     = var.object-key
  content = var.object-content
}
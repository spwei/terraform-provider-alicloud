variable "bucket-new" {
  default = "bucket-20180423-1"
}

variable "bucket-attr" {
  default = "bucket-20180423-2"
}

variable "acl-bj" {
  default = "public-read"
}

variable "index-doc" {
  default = "index.html"
}

variable "error-doc" {
  default = "error.html"
}

variable "target-prefix" {
  default = "log/"
}

variable "role-days" {
  default = "expirationByDays"
}

variable "rule-days" {
  default = 365
}

variable "role-date" {
  default = "expirationByDate"
}

variable "rule-date" {
  default = "2018-01-01"
}

variable "rule-prefix" {
  default = "path"
}

variable "allow-empty" {
  default = true
}

variable "referers" {
  default = "http://www.aliyun.com, https://www.aliyun.com, http://?.aliyun.com"
}

variable "allow-origins-star" {
  default = "*"
}

variable "allow-origins-aliyun" {
  default = "http://www.aliyun.com, http://*.aliyun.com"
}

variable "allow-methods-get" {
  default = "GET"
}

variable "allow-methods-put" {
  default = "PUT,GET"
}

variable "allowed_headers" {
  default = "authorization"
}

variable "expose_headers" {
  default = "x-oss-test, x-oss-test1"
}

variable "max_age_seconds" {
  default = 100
}

variable "object-key" {
  default = "object-content-key"
}

variable "object-content" {
  default = "This is my object content in May 22, 2017"
}
# Instance
variable "name" {
  description = "Name of your Kafka instance. The length cannot exceed 64 characters."
  default     = "terraform-example"
}

variable "topic_quota" {
  description = "The max num of topic can be create of the instance. When modify this value, it only adjust to a greater value."
  default     = 50
}

variable "disk_type" {
  description = "The disk type of the instance. 0: efficient cloud disk, 1: SSD."
  default     = 1
}

variable "disk_size" {
  description = "The disk size of the instance. When modify this value, it only adjust to a greater value."
  default     = 500
}

variable "deploy_type" {
  description = "The deploy type of the instance. Currently only support two deploy type, 4: eip/vpc instance, 5: vpc instance."
  default     = 5
}

variable "io_max" {
  description = "The peak value of io of the instance. When modify this value, it only support adjust to a greater value."
  default     = 20
}

variable "eip_max" {
  description = "The peak bandwidth of the instance. When modify this value, it only support adjust to a greater value."
  default     = 10
}


variable "paid_type" {
  description = "The paid type of the instance. Support two type, \"0\": pre paid type instance, \"1\": post paid type instance. When modify this value, it only support adjust from post pay to pre pay."
  default     = "PostPaid"
}

variable "spec_type" {
  description = "The spec type of the instance. Support two type, \"normal\": normal version instance, \"professional\": professional version instance. When modify this value, it only support adjust from normal to professional. Note only pre paid type instance support professional specific type."
  default     = "professional"
}

variable "vswitch_id" {
  description = "The ID of attaching vswitch to instance."
  default     = ""
}

# consumer group
variable "consumer_id" {
  description = "Id of ALIKAFKA consumer group. The length can't exceed 64 characters."
  default     = "terraform-example"
}

# sasl acl
variable "acl_username" {
  description = "Username of ALIKAFKA sasl acl. The length should between 1 to 64 characters."
  default     = "terraform"
}

variable "acl_resource_type" {
  description = "Resource type of ALIKAFKA sasl acl. The resource type can only be \"Topic\" and \"Group\"."
  default     = "Topic"
}

variable "acl_resource_name" {
  description = "Resource name of ALIKAFKA sasl acl. The resource name should be a topic or consumer group name."
  default     = "terraform-example"
}

variable "acl_resource_pattern_type" {
  description = "Resource pattern type of ALIKAFKA sasl acl. The resource pattern support two types \"LITERAL\" and \"PREFIXED\". \"LITERAL\": A literal name defines the full name of a resource. The special wildcard character \"*\" can be used to represent a resource with any name. \"PREFIXED\": A prefixed name defines a prefix for a resource."
  default     = "LITERAL"
}

variable "acl_operation_type" {
  description = "Acl operation type of ALIKAFKA sasl acl. The operation type can only be \"Write\" and \"Read\"."
  default     = "Write"
}
# sasl user
variable "username" {
  description = "Username of ALIKAFKA sasl user. The length should between 1 to 64 characters."
  default     = "terraform"
}

variable "password" {
  description = "Password of ALIKAFKA sasl user. The length should between 1 to 64 characters."
  default     = "YourPassword123"
}

# sasl topic
variable "topic" {
  description = "Name of ALIKAFKA topic. Two topics on a single instance cannot have the same name. The length cannot exceed 64 characters."
  default     = "terrform-example"
}

variable "local_topic" {
  description = "Whether the topic is localTopic or not."
  default     = false
}

variable "compact_topic" {
  description = "Whether the topic is compactTopic or not. Compact topic must be a localTopic."
  default     = false
}

variable "partition_num" {
  description = "The number of partitions of the topic. The number should between 1 and 48."
  default     = 6
}

variable "remark" {
  description = "This attribute is a concise description of topic. The length cannot exceed 64 characters."
  default     = "tf-example-alikafka-topic-remark"
}
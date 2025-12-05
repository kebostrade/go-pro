variable "cluster_name" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "cluster_version" {
  type    = string
  default = "1.28"
}

variable "node_groups" {
  type    = any
  default = {}
}

variable "enable_irsa" {
  type    = bool
  default = true
}

variable "environment" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}

variable "kms_key_arn" {
  type        = string
  description = "KMS key ARN for EKS secrets encryption. Required for production."
  default     = ""
}

variable "enable_secrets_encryption" {
  type        = bool
  description = "Enable secrets encryption at rest using KMS"
  default     = true
}

variable "enable_public_access" {
  type        = bool
  description = "Enable public API endpoint access (disable for production)"
  default     = false
}

variable "public_access_cidrs" {
  type        = list(string)
  description = "CIDR blocks allowed to access the public API endpoint"
  default     = []
}

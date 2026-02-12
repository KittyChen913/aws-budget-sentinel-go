variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "AWS CLI profile to use (leave null to use environment credentials, e.g. GitHub Actions)"
  type        = string
  nullable    = true
  default     = null
}

variable "function_name" {
  description = "Name of the Lambda function"
  type        = string
  default     = "aws-budget-sentinel"
}

variable "timeout" {
  description = "Lambda function timeout in seconds"
  type        = number
  default     = 30
}

variable "memory_size" {
  description = "Lambda function memory size in MB"
  type        = number
  default     = 128
}

variable "discord_webhook_url" {
  description = "Discord webhook URL for notifications (optional)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "log_retention_days" {
  description = "CloudWatch Logs retention period in days"
  type        = number
  default     = 7
}

variable "enable_schedule" {
  description = "Enable scheduled execution of Lambda function"
  type        = bool
  default     = true
}

variable "schedule_expression" {
  description = "CloudWatch Events schedule expression (e.g., rate(1 hour) or cron(0 * * * ? *))"
  type        = string
  default     = "rate(1 hour)"
}

variable "is_windows" {
  description = "Set to true if running on Windows"
  type        = bool
  default     = true
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default = {
    Project     = "aws-budget-sentinel"
    ManagedBy   = "Terraform"
    Environment = "production"
  }
}

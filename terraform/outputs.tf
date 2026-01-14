output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.budget_sentinel.arn
}

output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.budget_sentinel.function_name
}

output "lambda_role_arn" {
  description = "ARN of the Lambda execution role"
  value       = aws_iam_role.lambda_role.arn
}

output "log_group_name" {
  description = "Name of the CloudWatch Log Group"
  value       = aws_cloudwatch_log_group.lambda_logs.name
}

output "schedule_expression" {
  description = "Schedule expression for Lambda execution"
  value       = var.enable_schedule ? var.schedule_expression : "Not scheduled"
}

# CloudWatch Logs Group
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

# Lambda 函數
resource "aws_lambda_function" "budget_sentinel" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = var.function_name
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  runtime          = "provided.al2023"
  architectures    = ["x86_64"]
  timeout          = var.timeout
  memory_size      = var.memory_size

  environment {
    variables = merge(
      {
        # 可在此添加其他環境變數
      },
      var.discord_webhook_url != "" ? { DISCORD_WEBHOOK_URL = var.discord_webhook_url } : {}
    )
  }

  tags = var.tags
}

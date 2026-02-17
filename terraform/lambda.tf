# CloudWatch Logs Group
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

# Lambda 函數
# 注意：lambda-deployment.zip 由 CI/CD Build Job 編譯打包並下載到 ./lambda-artifact
resource "aws_lambda_function" "budget_sentinel" {
  filename         = "${path.module}/../lambda-artifact/lambda-deployment.zip"
  function_name    = var.function_name
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  source_code_hash = filebase64sha256("${path.module}/../lambda-artifact/lambda-deployment.zip")
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

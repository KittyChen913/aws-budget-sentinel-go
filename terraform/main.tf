terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile
}

# 編譯 Go Lambda 函數
resource "null_resource" "build_lambda" {
  triggers = {
    # 當程式碼變更時重新編譯
    go_mod      = filemd5("${path.module}/../go.mod")
    go_sum      = fileexists("${path.module}/../go.sum") ? filemd5("${path.module}/../go.sum") : ""
    main_go     = filemd5("${path.module}/../cmd/lambda/main.go")
    checks_dir  = md5(join("", [for f in fileset("${path.module}/../internal/checks", "**/*.go") : filemd5("${path.module}/../internal/checks/${f}")]))
    discord_dir = md5(join("", [for f in fileset("${path.module}/../internal/discord", "**/*.go") : filemd5("${path.module}/../internal/discord/${f}")]))
  }

  provisioner "local-exec" {
    command     = var.is_windows ? "powershell -File ${path.module}/build.ps1" : "bash ${path.module}/build.sh"
    working_dir = path.module
  }
}

# 打包 Lambda 部署檔案
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/../bootstrap"
  output_path = "${path.module}/lambda-deployment.zip"

  depends_on = [null_resource.build_lambda]
}

# IAM 角色
resource "aws_iam_role" "lambda_role" {
  name = "${var.function_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

# CloudWatch Logs 權限
resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# EC2 讀取權限
resource "aws_iam_role_policy" "lambda_ec2_read" {
  name = "${var.function_name}-ec2-read"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstances",
          "ec2:DescribeInstanceStatus"
        ]
        Resource = "*"
      }
    ]
  })
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

# CloudWatch Logs Group
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

# EventBridge (CloudWatch Events) 規則 - 定期執行
resource "aws_cloudwatch_event_rule" "schedule" {
  count               = var.enable_schedule ? 1 : 0
  name                = "${var.function_name}-schedule"
  description         = "Trigger ${var.function_name} on a schedule"
  schedule_expression = var.schedule_expression

  tags = var.tags
}

resource "aws_cloudwatch_event_target" "lambda_target" {
  count = var.enable_schedule ? 1 : 0
  rule  = aws_cloudwatch_event_rule.schedule[0].name
  arn   = aws_lambda_function.budget_sentinel.arn
}

resource "aws_lambda_permission" "allow_eventbridge" {
  count         = var.enable_schedule ? 1 : 0
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.budget_sentinel.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.schedule[0].arn
}

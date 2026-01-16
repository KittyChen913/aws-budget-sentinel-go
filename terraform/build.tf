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

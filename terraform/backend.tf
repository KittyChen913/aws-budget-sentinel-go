terraform {
  backend "s3" {
    bucket       = "kittychen913-terraform-state"
    key          = "aws-budget-sentinel/terraform.tfstate"
    region       = "us-east-1"
    use_lockfile = true
    encrypt      = true
  }
}

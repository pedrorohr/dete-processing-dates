terraform {
  required_version = "~> 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }

  backend "s3" {
    bucket = "pedrorohr-terraform-state"
    key    = "dete-processing-dates/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

data "aws_caller_identity" "current" {}

data "archive_file" "scraper_zip" {
  type        = "zip"
  source_file = "../bin/scraper"
  output_path = "bin/scraper.zip"
}

locals {
  account_id     = data.aws_caller_identity.current.account_id
  dete_processing_dates_url       = "https://enterprise.gov.ie/en/What-We-Do/Workplace-and-Skills/Employment-Permits/Current-Application-Processing-Dates/"
  lambda_handler = "scraper"
  name           = "dete-processing-dates"
  region         = "us-east-1"
}

data "aws_iam_policy_document" "assume_role" {
  policy_id = "${local.name}-lambda"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name                = "${local.name}-lambda"
  assume_role_policy  = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "logs" {
  policy_id = "${local.name}-lambda-logs"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = ["logs:CreateLogStream", "logs:PutLogEvents"]

    resources = [
      "arn:aws:logs:${local.region}:${local.account_id}:log-group:/aws/lambda/${local.name}*:*"
    ]
  }
}

resource "aws_iam_policy" "logs" {
  name   = "${local.name}-lambda-logs"
  policy = data.aws_iam_policy_document.logs.json
}

resource "aws_iam_role_policy_attachment" "logs" {
  depends_on = [aws_iam_role.lambda, aws_iam_policy.logs]
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.logs.arn
}

resource "aws_cloudwatch_log_group" "log" {
  name              = "/aws/lambda/${local.name}"
  retention_in_days = 7
}

resource "aws_lambda_function" "scraper" {
  filename          = data.archive_file.scraper_zip.output_path
  function_name     = local.name
  role              = aws_iam_role.lambda.arn
  handler           = local.lambda_handler
  source_code_hash  = filebase64sha256(data.archive_file.scraper_zip.output_path)
  runtime           = "go1.x"
  memory_size       = 1024
  timeout           = 30

  environment {
    variables = {
      DETE_PROCESSING_DATES_URL = local.dete_processing_dates_url
    }
  }
}

resource "aws_cloudwatch_event_rule" "every_five_minutes" {
    name = "every-five-minutes"
    description = "Fires every five minutes"
    schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "scrap_every_five_minutes" {
    rule = "${aws_cloudwatch_event_rule.every_five_minutes.name}"
    target_id = "scraper"
    arn = "${aws_lambda_function.scraper.arn}"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_scraper" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.scraper.function_name}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.every_five_minutes.arn}"
}

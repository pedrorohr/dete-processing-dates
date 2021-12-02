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
  account_id                = data.aws_caller_identity.current.account_id
  dete_processing_dates_url = "https://enterprise.gov.ie/en/What-We-Do/Workplace-and-Skills/Employment-Permits/Current-Application-Processing-Dates/"
  dete_bot_api_token        = "12345"
  dete_chat_id              = "1234"
  scraper_lambda_handler    = "scraper"
  notifier_lambda_handler   = "notifier"
  scraper_name              = "dete-processing-dates-scraper"
  notifier_name             = "dete-processing-dates-notifier"
  region                    = "us-east-1"
}

resource "aws_dynamodb_table" "dete-processing-dates" {
  name             = "DeteProcessingDates"
  read_capacity    = 2
  write_capacity   = 2
  hash_key         = "Type"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  attribute {
      name = "Type"
      type = "S"
  }
}

data "aws_iam_policy_document" "assume_role" {
  policy_id = "${local.scraper_name}-lambda"
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
  name               = "${local.scraper_name}-lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "logs" {
  policy_id = "${local.scraper_name}-lambda-logs"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = ["logs:CreateLogStream", "logs:PutLogEvents"]

    resources = [
      "arn:aws:logs:${local.region}:${local.account_id}:log-group:/aws/lambda/${local.scraper_name}*:*"
    ]
  }
}

resource "aws_iam_policy" "logs" {
  name   = "${local.scraper_name}-lambda-logs"
  policy = data.aws_iam_policy_document.logs.json
}

resource "aws_iam_role_policy_attachment" "logs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.logs.arn
}

data "aws_iam_policy_document" "dynamo" {
  policy_id = "${local.scraper_name}-lambda-dynamo"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem"
    ]

    resources = [aws_dynamodb_table.dete-processing-dates.arn]
  }
}

resource "aws_iam_policy" "dynamo" {
  name   = "${local.scraper_name}-lambda-dynamo"
  policy = data.aws_iam_policy_document.dynamo.json
}

resource "aws_iam_role_policy_attachment" "dynamo" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.dynamo.arn
}

resource "aws_cloudwatch_log_group" "log" {
  name              = "/aws/lambda/${local.scraper_name}"
  retention_in_days = 7
}

resource "aws_lambda_function" "scraper" {
  filename         = data.archive_file.scraper_zip.output_path
  function_name    = local.scraper_name
  role             = aws_iam_role.lambda.arn
  handler          = local.scraper_lambda_handler
  source_code_hash = filebase64sha256(data.archive_file.scraper_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30

  environment {
    variables = {
      DETE_PROCESSING_DATES_URL = local.dete_processing_dates_url
    }
  }
}

resource "aws_cloudwatch_event_rule" "every_five_minutes" {
    name                = "every-five-minutes"
    description         = "Fires every five minutes"
    schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "scrap_every_five_minutes" {
    rule      = aws_cloudwatch_event_rule.every_five_minutes.name
    target_id = "scraper"
    arn       = aws_lambda_function.scraper.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_scraper" {
    statement_id  = "AllowExecutionFromCloudWatch"
    action        = "lambda:InvokeFunction"
    function_name = aws_lambda_function.scraper.function_name
    principal     = "events.amazonaws.com"
    source_arn    = aws_cloudwatch_event_rule.every_five_minutes.arn
}

data "archive_file" "notifier_zip" {
  type        = "zip"
  source_file = "../bin/notifier"
  output_path = "bin/notifier.zip"
}

resource "aws_iam_role" "notifier-lambda" {
  name               = "${local.notifier_name}-lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "notifier-logs" {
  policy_id = "${local.notifier_name}-lambda-logs"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = ["logs:CreateLogStream", "logs:PutLogEvents"]

    resources = [
      "arn:aws:logs:${local.region}:${local.account_id}:log-group:/aws/lambda/${local.notifier_name}*:*"
    ]
  }
}

resource "aws_iam_policy" "notifier-logs" {
  name   = "${local.notifier_name}-lambda-logs"
  policy = data.aws_iam_policy_document.notifier-logs.json
}

resource "aws_iam_role_policy_attachment" "notifier-logs" {
  role       = aws_iam_role.notifier-lambda.name
  policy_arn = aws_iam_policy.notifier-logs.arn
}

data "aws_iam_policy_document" "notifier-dynamo" {
  policy_id = "${local.notifier_name}-lambda-dynamo"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = [
      "dynamodb:GetRecords",
      "dynamodb:GetShardIterator",
      "dynamodb:DescribeStream",
      "dynamodb:ListStreams"
    ]

    resources = [aws_dynamodb_table.dete-processing-dates.stream_arn]
  }
}

resource "aws_iam_policy" "notifier-dynamo" {
  name   = "${local.notifier_name}-lambda-dynamo"
  policy = data.aws_iam_policy_document.notifier-dynamo.json
}

resource "aws_iam_role_policy_attachment" "notifier-dynamo" {
  role       = aws_iam_role.notifier-lambda.name
  policy_arn = aws_iam_policy.notifier-dynamo.arn
}

resource "aws_cloudwatch_log_group" "notifier-log" {
  name              = "/aws/lambda/${local.notifier_name}"
  retention_in_days = 7
}

resource "aws_lambda_function" "notifier" {
  filename         = data.archive_file.notifier_zip.output_path
  function_name    = local.notifier_name
  role             = aws_iam_role.notifier-lambda.arn
  handler          = local.notifier_lambda_handler
  source_code_hash = filebase64sha256(data.archive_file.notifier_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30

  environment {
    variables = {
      DETE_BOT_API_TOKEN = local.dete_bot_api_token,
      DETE_CHAT_ID       = local.dete_chat_id
    }
  }
}

resource "aws_lambda_event_source_mapping" "notifier" {
  event_source_arn  = aws_dynamodb_table.dete-processing-dates.stream_arn
  function_name     = aws_lambda_function.notifier.arn
  starting_position = "LATEST"
}

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
  
  default_tags {
    tags = {
      Environment = var.environment
      Team        = "Pedro"
      Service     = "dete-processing-dates"
    }
  }
}

locals {
  scraper_service_name = "dete-processing-dates-scraper"
  notifier_service_name = "dete-processing-dates-notifier"
}

resource "aws_dynamodb_table" "dete_processing_dates" {
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

data "aws_iam_policy_document" "dynamo_table" {
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

    resources = [aws_dynamodb_table.dete_processing_dates.arn]
  }
}

data "aws_iam_policy_document" "dynamo_stream" {
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = [
      "dynamodb:GetRecords",
      "dynamodb:GetShardIterator",
      "dynamodb:DescribeStream",
      "dynamodb:ListStreams"
    ]

    resources = [aws_dynamodb_table.dete_processing_dates.stream_arn]
  }
}

module "scraper"{
  source         = "./modules/lambda"
  name           = local.scraper_service_name
  handler        = "scraper"
  source_file    = "../bin/scraper"
  extra_policies = {
    dynamo-table = data.aws_iam_policy_document.dynamo_table.json
  }
  env = {
    DETE_PROCESSING_DATES_URL = var.dete_processing_dates_url
  }
}

module "notifier"{
  source         = "./modules/lambda"
  name           = local.notifier_service_name
  handler        = "notifier"
  source_file    = "../bin/notifier"
  extra_policies = {
    dynamo-table = data.aws_iam_policy_document.dynamo_stream.json
  }
  env = {
    DETE_BOT_API_TOKEN = var.dete_bot_api_token,
    DETE_CHAT_ID       = var.dete_chat_id
  }
}

resource "aws_cloudwatch_event_rule" "every_five_minutes" {
  name                = "every-five-minutes"
  description         = "Fires every five minutes"
  schedule_expression = "rate(5 minutes)"

  tags = {
    Service = local.scraper_service_name
  }
}

resource "aws_cloudwatch_event_target" "scrap_every_five_minutes" {
  rule      = aws_cloudwatch_event_rule.every_five_minutes.name
  target_id = "scraper"
  arn       = module.scraper.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_scraper" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = module.scraper.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_five_minutes.arn
}

resource "aws_lambda_event_source_mapping" "notifier" {
  event_source_arn  = aws_dynamodb_table.dete_processing_dates.stream_arn
  function_name     = module.notifier.arn
  starting_position = "LATEST"
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

locals {
  account_id  = data.aws_caller_identity.current.account_id
  region      = data.aws_region.current.name
  output_path = "bin/${var.name}.zip"
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = var.source_file
  output_path = local.output_path
}

data "aws_iam_policy_document" "assume_role" {
  policy_id = "${var.name}-lambda"
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
  name               = "${var.name}-lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "logs" {
  policy_id = "${var.name}-lambda-logs"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = ["logs:CreateLogStream", "logs:PutLogEvents"]

    resources = [
      "arn:aws:logs:${local.region}:${local.account_id}:log-group:/aws/lambda/${var.name}*:*"
    ]
  }
}

resource "aws_iam_policy" "logs" {
  name   = "${var.name}-lambda-logs"
  policy = data.aws_iam_policy_document.logs.json
}

resource "aws_iam_role_policy_attachment" "logs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.logs.arn
}

resource "aws_iam_policy" "extra" {
  for_each = var.extra_policies
  name     = "${var.name}-lambda-${each.key}"
  policy   = each.value
}

resource "aws_iam_role_policy_attachment" "extra" {
  for_each   = var.extra_policies
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.extra[each.key].arn
}

resource "aws_cloudwatch_log_group" "log" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = 7
}

resource "aws_lambda_function" "lambda" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = var.name
  role             = aws_iam_role.lambda.arn
  handler          = var.handler
  source_code_hash = filebase64sha256(data.archive_file.lambda_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30

  environment {
    variables = var.env
  }
}

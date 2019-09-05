provider "archive" {}

resource "aws_s3_bucket" "a" {
  bucket = "${var.bucket}"
  acl    = "private"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "../build/main"
  output_path = "../build/main.zip"
}

data "aws_iam_policy_document" "policy" {
  statement {
    sid    = ""
    effect = "Allow"

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }

    actions = ["sts:AssumeRole"]
  }
}

# ROLES
# IAM role which dictates what other AWS services the Lambda fuc may access
resource "aws_iam_role" "lambda_role" {
  name               = "lambda_role"
  assume_role_policy = "${data.aws_iam_policy_document.policy.json}"
}

# POLICIES
# policy to be added to role, which allows to access resource
resource "aws_iam_role_policy" "dynamo_lambda_policy" {
  name = "dynamo_lambda_policy"
  role = "${aws_iam_role.lambda_role.id}"

  policy = <<EOF
{
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect": "Allow",
        "Action": ["dynamodb:*"],
        "Resource": "${var.dynamo_arn}"
      }
    ] 
}
EOF
}

resource "aws_iam_role_policy" "logs" {
  name = "${aws_iam_role.lambda_role.name}-logs"
  role = "${aws_iam_role.lambda_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect":"Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams",
        "logs:PutLogEvents",
        "logs:GetLogEvents",
        "logs:FilterLogEvents"
      ],
      "Resource":"*"
    }
  ]
}

  EOF
}

resource "aws_lambda_function" "lambda" {
  function_name = "images"
  filename      = "${data.archive_file.zip.output_path}"
  role          = "${aws_iam_role.lambda_role.arn}"
  handler       = "main"
  runtime       = "go1.x"

  environment {
    variables = {
      mode = "prod"
    }
  }
}

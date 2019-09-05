terraform {
  required_version = "~>0.11"
}

provider "aws" {
  region  = "${var.region}"
  version = "~> 1.0"
}

######################
# DynamoDB
#####################
module "dynamo" {
  source = "./modules/dynamo"
}

######################
# Lambda
######################
module "lambda" {
  source     = "./modules/lambda"
  bucket     = "${var.bucket}"
  dynamo_arn = "${module.dynamo.arn}"
}

######################
# API Gateway
######################
module "api" {
  source            = "./modules/apigw"
  lambda_arn        = "${module.lambda.arn}"
  lambda_invoke_arn = "${module.lambda.invoke_arn}"
}

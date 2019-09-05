# API Gateway
resource "aws_api_gateway_rest_api" "api" {
  name               = "ServerlessExample"
  binary_media_types = ["application/protobuf"]
  description        = "Terraform serverless"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  parent_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  path_part   = "{proxy+}"
}

# Request mapping
resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_resource.proxy.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "request_integration" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  resource_id = "${aws_api_gateway_method.proxy.resource_id}"
  http_method = "${aws_api_gateway_method.proxy.http_method}"

  type = "AWS_PROXY"

  # uri                     = "${aws_lambda_function.lambda.invoke_arn}"
  uri                     = "${var.lambda_invoke_arn}"
  integration_http_method = "POST"
}

resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  stage_name  = "test"

  depends_on = [
    "aws_api_gateway_integration.request_integration",
  ]
}

resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id = "AllowExecutionFromApiGateway"
  action       = "lambda:InvokeFunction"

  # function_name = "${aws_lambda_function.lambda.function_name}"
  function_name = "${var.lambda_arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"

  depends_on = [
    "aws_api_gateway_rest_api.api",
  ]
}

output "base_url" {
  value = "${aws_api_gateway_deployment.deployment.invoke_url}"
}

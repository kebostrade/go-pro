# Terraform configuration for AWS infrastructure

terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "terraform-state-bucket"
    key    = "cloud-cicd-go/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.region
}

# Variables
variable "region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "app_name" {
  description = "Application name"
  type        = string
  default     = "cloud-cicd-app"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

# S3 Bucket
resource "aws_s3_bucket" "app_bucket" {
  bucket = "${var.app_name}-${var.environment}-${data.aws_caller_identity.current.account_id}"

  tags = {
    Name        = "${var.app_name}-bucket"
    Environment = var.environment
  }
}

resource "aws_s3_bucket_versioning" "app_bucket_versioning" {
  bucket = aws_s3_bucket.app_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "app_bucket_lifecycle" {
  bucket = aws_s3_bucket.app_bucket.id

  rule {
    id     = "delete-old-versions"
    status = "Enabled"

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }
}

# DynamoDB Table
resource "aws_dynamodb_table" "app_table" {
  name         = "${var.app_name}-${var.environment}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "email"
    type = "S"
  }

  global_secondary_index {
    name            = "EmailIndex"
    hash_key        = "email"
    projection_type = "ALL"
  }

  # Point-in-Time Recovery: automatic backup for disaster recovery
  point_in_time_recovery_specification {
    enabled = true
  }

  tags = {
    Name        = "${var.app_name}-table"
    Environment = var.environment
  }
}

# SQS Queue
resource "aws_sqs_queue" "app_queue" {
  name                      = "${var.app_name}-${var.environment}-queue"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 345600
  receive_wait_time_seconds = 10

  tags = {
    Name        = "${var.app_name}-queue"
    Environment = var.environment
  }
}

# SQS Dead Letter Queue
resource "aws_sqs_queue" "app_dlq" {
  name = "${var.app_name}-${var.environment}-dlq"

  tags = {
    Name        = "${var.app_name}-dlq"
    Environment = var.environment
  }
}

resource "aws_sqs_queue_redrive_policy" "app_queue_redrive" {
  queue_url = aws_sqs_queue.app_queue.id

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.app_dlq.arn
    maxReceiveCount     = 3
  })
}

# SNS Topic
resource "aws_sns_topic" "app_topic" {
  name              = "${var.app_name}-${var.environment}-topic"
  kms_master_key_id = "alias/aws/sns"

  tags = {
    Name        = "${var.app_name}-topic"
    Environment = var.environment
  }
}

# SNS Subscription to SQS
resource "aws_sns_topic_subscription" "app_topic_sqs" {
  topic_arn = aws_sns_topic.app_topic.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.app_queue.arn
}

# IAM Role for Lambda
resource "aws_iam_role" "lambda_role" {
  name = "${var.app_name}-lambda-role"

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

  tags = {
    Name        = "${var.app_name}-lambda-role"
    Environment = var.environment
  }
}

# IAM Policy for Lambda
resource "aws_iam_role_policy" "lambda_policy" {
  name = "${var.app_name}-lambda-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.app_bucket.arn}/*"
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = aws_dynamodb_table.app_table.arn
      },
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = aws_sqs_queue.app_queue.arn
      },
      {
        Effect = "Allow"
        Action = [
          "sns:Publish"
        ]
        Resource = aws_sns_topic.app_topic.arn
      }
    ]
  })
}

# Lambda Function
resource "aws_lambda_function" "app_function" {
  filename      = "../../aws/lambda/function.zip"
  function_name = "${var.app_name}-${var.environment}"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2"

  source_code_hash = fileexists("../../aws/lambda/function.zip") ? filebase64sha256("../../aws/lambda/function.zip") : null

  environment {
    variables = {
      ENV                = var.environment
      AWS_BUCKET_NAME    = aws_s3_bucket.app_bucket.id
      AWS_DYNAMODB_TABLE = aws_dynamodb_table.app_table.name
      AWS_SQS_QUEUE_URL  = aws_sqs_queue.app_queue.url
      AWS_SNS_TOPIC_ARN  = aws_sns_topic.app_topic.arn
    }
  }

  timeout     = 30
  memory_size = 256

  tags = {
    Name        = "${var.app_name}-function"
    Environment = var.environment
  }
}

# Lambda Function URL
resource "aws_lambda_function_url" "app_function_url" {
  function_name      = aws_lambda_function.app_function.function_name
  authorization_type = "NONE"

  cors {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PUT", "DELETE"]
    allow_headers = ["*"]
    max_age       = 86400
  }
}

# API Gateway (optional)
resource "aws_apigatewayv2_api" "app_api" {
  name          = "${var.app_name}-${var.environment}-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["*"]
    max_age       = 300
  }

  tags = {
    Name        = "${var.app_name}-api"
    Environment = var.environment
  }
}

resource "aws_apigatewayv2_integration" "app_integration" {
  api_id           = aws_apigatewayv2_api.app_api.id
  integration_type = "AWS_PROXY"

  integration_uri    = aws_lambda_function.app_function.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "app_route" {
  api_id    = aws_apigatewayv2_api.app_api.id
  route_key = "$default"

  target = "integrations/${aws_apigatewayv2_integration.app_integration.id}"
}

resource "aws_apigatewayv2_stage" "app_stage" {
  api_id      = aws_apigatewayv2_api.app_api.id
  name        = var.environment
  auto_deploy = true

  tags = {
    Name        = "${var.app_name}-stage"
    Environment = var.environment
  }
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.app_function.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.app_api.execution_arn}/*/*"
}

# Data source for current AWS account
data "aws_caller_identity" "current" {}

# Outputs
output "lambda_function_url" {
  description = "Lambda function URL"
  value       = aws_lambda_function_url.app_function_url.function_url
}

output "api_gateway_url" {
  description = "API Gateway URL"
  value       = aws_apigatewayv2_stage.app_stage.invoke_url
}

output "s3_bucket_name" {
  description = "S3 bucket name"
  value       = aws_s3_bucket.app_bucket.id
}

output "dynamodb_table_name" {
  description = "DynamoDB table name"
  value       = aws_dynamodb_table.app_table.name
}

output "sqs_queue_url" {
  description = "SQS queue URL"
  value       = aws_sqs_queue.app_queue.url
}

output "sns_topic_arn" {
  description = "SNS topic ARN"
  value       = aws_sns_topic.app_topic.arn
}


# GO-PRO Learning Platform - Main Terraform Configuration
# This file orchestrates the infrastructure deployment across multiple cloud providers

terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 7.12"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }

  # Backend configuration for state management
  # Uncomment and configure when ready to use remote state
  # backend "s3" {
  #   bucket         = "gopro-terraform-state"
  #   key            = "infrastructure/terraform.tfstate"
  #   region         = "us-east-1"
  #   encrypt        = true
  #   dynamodb_table = "gopro-terraform-locks"
  # }
}

# Provider configurations
provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "GO-PRO"
      Environment = var.environment
      ManagedBy   = "Terraform"
      Owner       = var.owner
    }
  }
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# Data sources
data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}

# Local variables
locals {
  name_prefix = "gopro-${var.environment}"

  common_tags = {
    Project     = "GO-PRO"
    Environment = var.environment
    ManagedBy   = "Terraform"
    Owner       = var.owner
    CostCenter  = var.cost_center
  }

  # Network configuration
  vpc_cidr = var.vpc_cidr
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  # Database configuration
  db_name     = "gopro_${var.environment}"
  db_username = "gopro_admin"

  # Redis configuration
  redis_node_type = var.environment == "production" ? "cache.r6g.large" : "cache.t4g.micro"

  # Kubernetes configuration
  cluster_name = "${local.name_prefix}-eks"
}

# Random password generation
resource "random_password" "db_password" {
  length  = 32
  special = true
}

resource "random_password" "redis_auth_token" {
  length  = 32
  special = false
}

# VPC Module
module "vpc" {
  source = "./modules/vpc"

  name_prefix = local.name_prefix
  vpc_cidr    = local.vpc_cidr
  azs         = local.azs
  environment = var.environment
  tags        = local.common_tags
}

# Security Groups Module
module "security_groups" {
  source = "./modules/security-groups"

  name_prefix = local.name_prefix
  vpc_id      = module.vpc.vpc_id
  environment = var.environment
  tags        = local.common_tags
}

# RDS PostgreSQL Module
module "rds" {
  source = "./modules/rds"

  name_prefix           = local.name_prefix
  vpc_id                = module.vpc.vpc_id
  subnet_ids            = module.vpc.private_subnet_ids
  security_group_ids    = [module.security_groups.rds_security_group_id]
  db_name               = local.db_name
  db_username           = local.db_username
  db_password           = random_password.db_password.result
  instance_class        = var.rds_instance_class
  allocated_storage     = var.rds_allocated_storage
  multi_az              = var.environment == "production"
  backup_retention_days = var.environment == "production" ? 30 : 7
  environment           = var.environment
  tags                  = local.common_tags
}

# ElastiCache Redis Module
module "redis" {
  source = "./modules/redis"

  name_prefix        = local.name_prefix
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnet_ids
  security_group_ids = [module.security_groups.redis_security_group_id]
  node_type          = local.redis_node_type
  num_cache_nodes    = var.environment == "production" ? 3 : 1
  auth_token         = random_password.redis_auth_token.result
  environment        = var.environment
  tags               = local.common_tags
}

# EKS Cluster Module
module "eks" {
  source = "./modules/eks"

  cluster_name    = local.cluster_name
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnet_ids
  cluster_version = var.eks_cluster_version
  node_groups     = var.eks_node_groups
  enable_irsa     = true
  environment     = var.environment
  tags            = local.common_tags
}

# MSK (Managed Kafka) Module
module "msk" {
  source = "./modules/msk"

  name_prefix        = local.name_prefix
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnet_ids
  security_group_ids = [module.security_groups.msk_security_group_id]
  kafka_version      = var.msk_kafka_version
  broker_node_count  = var.environment == "production" ? 3 : 2
  instance_type      = var.msk_instance_type
  ebs_volume_size    = var.msk_ebs_volume_size
  environment        = var.environment
  tags               = local.common_tags
}

# S3 Buckets Module
module "s3" {
  source = "./modules/s3"

  name_prefix = local.name_prefix
  environment = var.environment
  tags        = local.common_tags
}

# CloudWatch Module
module "cloudwatch" {
  source = "./modules/cloudwatch"

  name_prefix  = local.name_prefix
  cluster_name = module.eks.cluster_name
  environment  = var.environment
  tags         = local.common_tags
}

# Secrets Manager Module
module "secrets" {
  source = "./modules/secrets"

  name_prefix    = local.name_prefix
  db_password    = random_password.db_password.result
  redis_token    = random_password.redis_auth_token.result
  db_endpoint    = module.rds.endpoint
  redis_endpoint = module.redis.endpoint
  environment    = var.environment
  tags           = local.common_tags
}

# IAM Roles Module
module "iam" {
  source = "./modules/iam"

  name_prefix  = local.name_prefix
  cluster_name = module.eks.cluster_name
  environment  = var.environment
  tags         = local.common_tags
}

# Route53 DNS Module (if domain is provided)
module "route53" {
  source = "./modules/route53"
  count  = var.domain_name != "" ? 1 : 0

  domain_name = var.domain_name
  environment = var.environment
  tags        = local.common_tags
}

# ACM Certificate Module (if domain is provided)
module "acm" {
  source = "./modules/acm"
  count  = var.domain_name != "" ? 1 : 0

  domain_name = var.domain_name
  zone_id     = module.route53[0].zone_id
  environment = var.environment
  tags        = local.common_tags
}

# Application Load Balancer Module
module "alb" {
  source = "./modules/alb"

  name_prefix        = local.name_prefix
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.public_subnet_ids
  security_group_ids = [module.security_groups.alb_security_group_id]
  certificate_arn    = var.domain_name != "" ? module.acm[0].certificate_arn : null
  environment        = var.environment
  tags               = local.common_tags
}

# WAF Module (for production)
module "waf" {
  source = "./modules/waf"
  count  = var.environment == "production" ? 1 : 0

  name_prefix = local.name_prefix
  alb_arn     = module.alb.arn
  environment = var.environment
  tags        = local.common_tags
}

# Backup Module
module "backup" {
  source = "./modules/backup"
  count  = var.environment == "production" ? 1 : 0

  name_prefix = local.name_prefix
  rds_arn     = module.rds.arn
  environment = var.environment
  tags        = local.common_tags
}

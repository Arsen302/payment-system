provider "aws" {
  region = var.aws_region
}

terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11"
    }
  }

  backend "s3" {
    bucket         = "payment-system-terraform-state-dev"
    key            = "terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "payment-system-terraform-lock-dev"
  }
}

module "vpc" {
  source = "../../modules/vpc"

  environment       = var.environment
  vpc_cidr          = var.vpc_cidr
  availability_zones = var.availability_zones
  public_subnets    = var.public_subnets
  private_subnets   = var.private_subnets
}

module "eks" {
  source = "../../modules/eks"

  environment         = var.environment
  cluster_name        = var.cluster_name
  kubernetes_version  = var.kubernetes_version
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  node_instance_types = var.node_instance_types
  node_desired_size   = var.node_desired_size
  node_min_size       = var.node_min_size
  node_max_size       = var.node_max_size
}

module "rds" {
  source = "../../modules/rds"

  environment        = var.environment
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnet_ids
  instance_class     = var.db_instance_class
  allocated_storage  = var.db_allocated_storage
  engine_version     = var.db_engine_version
  database_name      = var.database_name
  master_username    = var.db_master_username
}

module "elasticache" {
  source = "../../modules/elasticache"

  environment       = var.environment
  vpc_id            = module.vpc.vpc_id
  subnet_ids        = module.vpc.private_subnet_ids
  node_type         = var.redis_node_type
  engine_version    = var.redis_engine_version
}

module "msk" {
  source = "../../modules/msk"

  environment       = var.environment
  vpc_id            = module.vpc.vpc_id
  subnet_ids        = module.vpc.private_subnet_ids
  kafka_version     = var.kafka_version
  broker_node_type  = var.kafka_broker_node_type
  broker_count      = var.kafka_broker_count
}

module "argocd" {
  source = "../../modules/argocd"

  environment       = var.environment
  eks_cluster_name  = module.eks.cluster_name
  chart_version     = var.argocd_chart_version
  admin_password    = var.argocd_admin_password
}

module "monitoring" {
  source = "../../modules/monitoring"

  environment       = var.environment
  eks_cluster_name  = module.eks.cluster_name
}

output "eks_cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "rds_endpoint" {
  value = module.rds.endpoint
}

output "elasticache_endpoint" {
  value = module.elasticache.endpoint
}

output "msk_bootstrap_brokers" {
  value = module.msk.bootstrap_brokers
} 
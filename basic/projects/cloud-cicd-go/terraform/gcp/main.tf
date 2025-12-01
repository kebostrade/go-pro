# Terraform configuration for GCP infrastructure

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
  
  backend "gcs" {
    bucket = "terraform-state-bucket"
    prefix = "cloud-cicd-go"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Variables
variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
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

# Enable required APIs
resource "google_project_service" "required_apis" {
  for_each = toset([
    "run.googleapis.com",
    "cloudfunctions.googleapis.com",
    "storage.googleapis.com",
    "pubsub.googleapis.com",
    "firestore.googleapis.com",
    "container.googleapis.com",
  ])
  
  service            = each.value
  disable_on_destroy = false
}

# Cloud Storage Bucket
resource "google_storage_bucket" "app_bucket" {
  name          = "${var.project_id}-${var.app_name}-${var.environment}"
  location      = var.region
  force_destroy = true
  
  uniform_bucket_level_access = true
  
  versioning {
    enabled = true
  }
  
  lifecycle_rule {
    condition {
      age = 30
    }
    action {
      type = "Delete"
    }
  }
  
  labels = {
    environment = var.environment
    app         = var.app_name
  }
}

# Pub/Sub Topic
resource "google_pubsub_topic" "app_topic" {
  name = "${var.app_name}-${var.environment}-topic"
  
  labels = {
    environment = var.environment
    app         = var.app_name
  }
}

# Pub/Sub Subscription
resource "google_pubsub_subscription" "app_subscription" {
  name  = "${var.app_name}-${var.environment}-subscription"
  topic = google_pubsub_topic.app_topic.name
  
  ack_deadline_seconds = 20
  
  retry_policy {
    minimum_backoff = "10s"
    maximum_backoff = "600s"
  }
  
  labels = {
    environment = var.environment
    app         = var.app_name
  }
}

# Firestore Database
resource "google_firestore_database" "app_database" {
  project                     = var.project_id
  name                        = "(default)"
  location_id                 = var.region
  type                        = "FIRESTORE_NATIVE"
  deletion_policy             = "DELETE"
  delete_protection_state     = "DELETE_PROTECTION_DISABLED"

  depends_on = [google_project_service.required_apis]
}

# Service Account for Cloud Run
resource "google_service_account" "cloud_run_sa" {
  account_id   = "${var.app_name}-cloud-run"
  display_name = "Cloud Run Service Account"
  description  = "Service account for Cloud Run application"
}

# IAM bindings for Service Account
resource "google_project_iam_member" "cloud_run_sa_roles" {
  for_each = toset([
    "roles/storage.objectViewer",
    "roles/pubsub.publisher",
    "roles/datastore.user",
  ])
  
  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}

# Cloud Run Service
resource "google_cloud_run_service" "app" {
  name     = var.app_name
  location = var.region
  
  template {
    spec {
      service_account_name = google_service_account.cloud_run_sa.email
      
      containers {
        image = "gcr.io/${var.project_id}/${var.app_name}:latest"
        
        ports {
          container_port = 8080
        }
        
        env {
          name  = "ENV"
          value = var.environment
        }
        
        env {
          name  = "GCP_PROJECT_ID"
          value = var.project_id
        }
        
        env {
          name  = "GCP_BUCKET_NAME"
          value = google_storage_bucket.app_bucket.name
        }
        
        env {
          name  = "GCP_PUBSUB_TOPIC"
          value = google_pubsub_topic.app_topic.name
        }
        
        resources {
          limits = {
            cpu    = "1000m"
            memory = "512Mi"
          }
        }
      }
    }
    
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "10"
        "autoscaling.knative.dev/minScale" = "1"
      }
    }
  }
  
  traffic {
    percent         = 100
    latest_revision = true
  }
  
  depends_on = [google_project_service.required_apis]
}

# Allow unauthenticated access to Cloud Run
resource "google_cloud_run_service_iam_member" "public_access" {
  service  = google_cloud_run_service.app.name
  location = google_cloud_run_service.app.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# GKE Cluster (optional)
resource "google_container_cluster" "primary" {
  name     = "${var.app_name}-gke-${var.environment}"
  location = var.region
  project  = var.project_id

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1

  network    = "default"
  subnetwork = "default"

  # VPC-native cluster with secondary IP ranges (modern GKE best practice)
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  # Disable legacy client certificate authentication for security
  # Client certificate auth is deprecated and increases attack surface
  master_auth {
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  # Logging and monitoring configuration
  logging_service    = "logging.googleapis.com/kubernetes"
  monitoring_service = "monitoring.googleapis.com/kubernetes"

  # Required for GKE - set to false for dev environments
  deletion_protection = false

  depends_on = [google_project_service.required_apis]
}

# GKE Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name       = "${var.app_name}-node-pool"
  location   = var.region
  project    = var.project_id
  cluster    = google_container_cluster.primary.name
  node_count = 1
  
  autoscaling {
    min_node_count = 1
    max_node_count = 3
  }
  
  node_config {
    preemptible  = true
    machine_type = "e2-medium"
    
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    
    labels = {
      environment = var.environment
      app         = var.app_name
    }
    
    tags = ["gke-node", var.app_name]
  }
}

# Outputs
output "cloud_run_url" {
  description = "Cloud Run service URL"
  value       = google_cloud_run_service.app.status[0].url
}

output "bucket_name" {
  description = "Cloud Storage bucket name"
  value       = google_storage_bucket.app_bucket.name
}

output "pubsub_topic" {
  description = "Pub/Sub topic name"
  value       = google_pubsub_topic.app_topic.name
}

output "gke_cluster_name" {
  description = "GKE cluster name"
  value       = google_container_cluster.primary.name
}

output "gke_cluster_endpoint" {
  description = "GKE cluster endpoint"
  value       = google_container_cluster.primary.endpoint
  sensitive   = true
}


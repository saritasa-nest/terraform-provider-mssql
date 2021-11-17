terraform {
  required_providers {
    mssql = {
      source  = "saritasa/mssql"
      version = "~> 0.1.0"
    }
  }
  required_version = ">= 0.13"
}

provider "mssql" {
  endpoint = var.host
  port     = var.port
  username = var.user
  password = var.password
}
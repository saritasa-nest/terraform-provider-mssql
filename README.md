# Microsoft SQL Server Terraform Provider

## Usage

```hcl
terraform {
  required_providers {
    mssql = {
      source  = "saritasa/provider"
      version = "~> 0.1.0"
    }
  }
  required_version = ">= 0.13"
}

provider "mssql" {
  endpoint = "localhost"
  username = "admin"
  password = "mypass"
}
```

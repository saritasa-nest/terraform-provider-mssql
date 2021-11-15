**This repository is an unofficial fork**

---

Terraform Provider
==================

Usage
-----

```hcl
terraform {
  required_providers {
    mysql = {
      source  = "saritasa/mssql"
      version = "~> 0.1.0"
    }
  }
  required_version = ">= 0.13"
}

provider "mssql" {
  endpoint = "localhost"
  username = "admin"
}
```

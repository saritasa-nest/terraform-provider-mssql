# Microsoft SQL Server Terraform Provider

## Example Usage

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

## Argument Reference

* **endpoint** - MSSQL server host. Can be set via environment variable `MSSQL_ENDPOINT`
* **port** - MSSQL server port. Default 1433. Can be set via environment variable `MSSQL_PORT`
* **username** - will be used to connect to MSSQL server. Default `sa`. Can be set via environment variable `MSSQL_USERNAME`
* **password** - will be used to connect to MSSQL server. Can be set via environment variable `MSSQL_PASSWORD`
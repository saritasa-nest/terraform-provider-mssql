---
layout: "mssql"
page_title: "MS SQL: mssql_login"
sidebar_current: "docs-mssql-resource-login"
description: |-
Creates and manages a login (user) in MS SQL server
---

# mssql\_login

The `mssql_database` resource creates and manages a database on a MS SQL
server.

```hql
resource "mssql_login" "demo" {
  name = "demo_login"
  password = "!12345678p"
  options = {
    default_database = "mydb"
  }

  depends_on = [mssql_database.mydb]
}

output "demo_login" {
  value = mssql_login.demo
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database. This must be unique within
  a given MS SQL server.
* `password` - (Required) password to set for user
* `options` - (Optional) - a key-value map of options supported by DB engine for logins
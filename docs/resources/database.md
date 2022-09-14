---
layout: "mssql"
page_title: "MS SQL: mssql_database"
sidebar_current: "docs-mssql-resource-database"
description: |-
  Creates and manages a database on a MS SQL server.
---

# mssql\_database

The ``mssql_database`` resource creates and manages a database on a MS SQL
server.

~> **Caution:** The ``mssql_database`` resource can completely delete your
database just as easily as it can create it. To avoid costly accidents,
consider setting
[``prevent_destroy``](/docs/configuration/resources.html#prevent_destroy)
on your database resources as an extra safety measure.

## Example Usage

```hcl
resource "mssql_database" "app" {
  name = "my_awesome_app"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database. This must be unique within
  a given MS SQL server and may or may not be case-sensitive depending on
  the operating system on which the MS SQL server is running.

* `default_collation` - (Optional) The default collation to use when a table
  is created without specifying an explicit collation. Defaults to
  ``Latin1_General_100_CI_AI``. Each character set has its own set of collations, so
  changing the character set requires also changing the collation.

* `options` - (Optional) a key-value map of options supported by MSSQL to pass to 
  database engine on DB creation. Ex. 
  ```hql
  resource "mssql_database" "mydb" {
    name = "mydb"
    default_collation = "Latin1_General_CI_AS"
  
    options = {
      default_language = "us_english"
    }
  }
  ```
  **Note:** This feature is incomplete. May have issues on state update. 

Note that the defaults for character set and collation above do not respect
any defaults set on the MS SQL server, so that the configuration can be set
appropriately even though Terraform cannot see the server-level defaults. If
you wish to use the server's defaults you must consult the server's
configuration and then set the ``default_character_set`` and
``default_collation`` to match.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the database.
* `id` - The id of the database.
* `default_character_set` - The default_character_set of the database.
* `default_collation` - The default_collation of the database.

## Import

Databases can be imported using their name, e.g.

```
$ terraform import mssql_database.example my-example-database
```

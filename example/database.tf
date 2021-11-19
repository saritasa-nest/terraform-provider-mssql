resource "mssql_database" "mydb" {
  name = "mydb"
}

output "mydb" {
  value = {
    name = mssql_database.mydb.name
    default_language = mssql_database.mydb.default_language
  }
}
resource "mssql_database" "mydb" {
  name = "mydb"
#  default_collation = "Latin1_General_CI_AS"
  options = {
#    default_language = "us_english"
  }
}

output "mydb" {
  value = mssql_database.mydb
}
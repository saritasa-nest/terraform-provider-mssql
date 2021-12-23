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

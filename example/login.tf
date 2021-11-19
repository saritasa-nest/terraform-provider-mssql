resource "mssql_login" "demo" {
  name = "demo_login"
  password = "123456"
  options = {
    default_database = "mydb"
  }
}

output "demo_login" {
  value = mssql_login.demo
}

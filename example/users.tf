resource "mssql_user" "user_global" {
  database = mssql_database.mydb.id

  username = "demo_user"
  login_name = "demo_login"
}

output "demo_user" {
  value = mssql_user.user_global
}

resource "mssql_user" "user_with_password" {
  database = mssql_database.mydb.id

  username = "demo_user1"
  password = "123456"
}

output "user_with_password" {
  value = mssql_user.user_with_password
}

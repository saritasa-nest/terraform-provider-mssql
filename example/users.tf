resource "mssql_user" "user_global" {
  database = mssql_database.mydb.id

  username = "demo_user"
  login_name = "demo_login"

  depends_on = [mssql_login.demo]
}

output "demo_user" {
  value = mssql_user.user_global
  sensitive = true
}

# Create user with password (without creating 'login' previously) works only if "CONTAINED DATABASE AUTHENTICATION == 1"
# EXEC sp_configure 'CONTAINED DATABASE AUTHENTICATION', 1
# to check:
# EXEC sp_configure 'CONTAINED DATABASE AUTHENTICATION'
# See: https://stackoverflow.com/questions/20030612/you-can-only-create-a-user-with-a-password-in-a-contained-database
resource "mssql_user" "user_with_password" {
  database = mssql_database.mydb.id

  username = "demo_user1"
  password = "123456"
}

output "user_with_password" {
  value = mssql_user.user_with_password
  sensitive = true
}

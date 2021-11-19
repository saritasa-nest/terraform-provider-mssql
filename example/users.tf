resource "mssql_user" "user" {
  database = mssql_database.mydb.id

  username = "demo_user"
  login_name = "demo_login"
}



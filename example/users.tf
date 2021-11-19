resource "mssql_user" "demo_user1" {
  username = "demo_user1"
  password = "123456"
}

resource "mssql_login" "demo_user2" {
  name = "demo_user1"
  password = "123456"
  default_database = "dbo"
}


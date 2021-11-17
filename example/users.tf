resource "mssql_user" "demo_user1" {
  username = "demo_user1"
  password = "123456"
}

resource "mssql_user" "demo_user2" {
  username = "demo_user1"
  password = "123456"
  default_schema = "dbo"
}

resource "mssql_user" "demo_guest" {
  username = "guest"
  password = "123456"
  default_schema = "dbo"
}
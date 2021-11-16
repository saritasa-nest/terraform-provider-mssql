variable "host" {
  type = string
  description = "MSSQL Database Host"
}

variable "port" {
  type = number
  description = "MSSQL Database Port"
  default = 1433
}

variable "user" {
  type = string
  description = "MSSQL Database User - on behalf of this user terraform will act"
  default = "sa"
}

variable "password" {
  type = string
  description = "MSSQL Database User Password - terraform will use this password to connect an manipulate DB"
}
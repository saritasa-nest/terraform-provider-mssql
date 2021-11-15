package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/saritasa/terraform-provider-mssql/mssql"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mssql.Provider})
}

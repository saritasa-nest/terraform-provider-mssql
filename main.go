package main

import (
	"context"
	"flag"
	"github.com/saritasa/terraform-provider-mssql/mssql"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: mssql.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/saritasa/mssql", opts)

		if err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(opts)
}

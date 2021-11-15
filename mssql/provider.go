package mssql

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSSQL_ENDPOINT", nil),
				Description: "MSSQL server host",
			},

			"port": {
				Type:        schema.TypeInt,
				Required:    false,
				DefaultFunc: schema.EnvDefaultFunc("MSSQL_PORT", 1433),
				Description: "MSSQL server port",
			},

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSSQL_USERNAME", nil),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSSQL_PASSWORD", nil),
			},

			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSSQL_DATABASE", nil),
			},

			"max_conn_lifetime_sec": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"max_open_conns": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"connect_retry_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"mysql_tables": DataSourceTables(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"mssql_database": ResourceDatabase(),
			"mssql_grant":    ResourceGrant(),
			"mssql_role":     ResourceRole(),
			"mssql_user":     ResourceUser(),
			"mssql_sql":      ResourceSql(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var endpoint = d.Get("endpoint").(string)

	client := MsSqlClient{
		Host:     endpoint,
		Port:     d.Get("port").(int),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Database: d.Get("database").(string),
	}

	return client, diag.Diagnostics{}
}

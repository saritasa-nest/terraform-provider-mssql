package provider

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/saritasa/terraform-provider-mssql/model"
	"github.com/saritasa/terraform-provider-mssql/mssql"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateDatabase,
		UpdateContext: UpdateDatabase,
		ReadContext:   ReadDatabase,
		DeleteContext: DeleteDatabase,
		Importer: &schema.ResourceImporter{
			StateContext: ImportDatabase,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"default_collation": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:  "SQL_Latin1_General_CP1_CI_AS",
			},

			"options": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func CreateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	database := new(model.Database).Parse(d)

	stmtSQL := fmt.Sprintf("CREATE DATABASE [%s]", database.Name)
	if database.DefaultCollation != "" {
		stmtSQL += " COLLATE " + database.DefaultCollation
	}
	if len(database.Options) > 0 {
		stmtSQL += " WITH "
		for opt := range database.Options {
			stmtSQL += fmt.Sprintf(" %s = %s,", opt, database.Options[opt].ValueOrSqlNull())
		}
		stmtSQL = strings.TrimRight(stmtSQL, ",")
	}

	err := connector.ExecContext(ctx, stmtSQL)
	if err == nil {
		d.SetId(database.Name)
	}
	return diag.FromErr(err)
}

func UpdateDatabase(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	diags := diag.Diagnostics{}

	database := new(model.Database).Parse(data)

	if data.HasChange("default_collation") && database.DefaultCollation != "" {
		smtSQL := fmt.Sprintf("ALTER DATABASE %s COLLATE %s", database.Name, database.DefaultCollation)
		err := connector.ExecContext(ctx, smtSQL)
		if err != nil {
			diags = diag.FromErr(err)
		}
	}

	if data.HasChange("options") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "MSSQL database updates may not function properly, especially options remove",
			Detail:   "MSSQL options update is not versatile, may not detect and apply correctly",
		})
		for opt := range database.Options {
			value := database.Options[opt].ValueOrSqlNull()
			stmtSQL := fmt.Sprintf("ALTER DATABASE [%s] WITH %s = %s", database.Name, opt, value)
			err := connector.ExecContext(ctx, stmtSQL)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Summary: fmt.Sprintf("MSSQL database %s option '%s' update", database.Name, opt),
					Detail:  err.Error(),
				})
			}
		}
	}

	return diags
}

func ReadDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	database := new(model.Database).Parse(d)
	dbName := d.Id()

	stmtSQL := "SELECT name, collation_name FROM sys.databases WHERE name LIKE '" + d.Id() + "'"

	log.Println("Executing statement:", stmtSQL)
	var collation model.NullString
	err := connector.QueryRowContext(ctx, stmtSQL, func(row *sql.Row) error {
		return row.Scan(&dbName, &collation)
	})
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Summary: fmt.Sprintf("read database %s info: %s", dbName, err.Error()),
		}}
	} else {
		err = d.Set("default_collation", collation.ToString())
	}

	if err != nil {
		return diag.FromErr(err)
	}
	return database.ToSchema(d)
}

func DeleteDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)

	stmtSQL := "DROP DATABASE " + d.Get("name").(string)
	log.Println("Executing statement:", stmtSQL)
	err := connector.ExecContext(ctx, stmtSQL)

	if err == nil {
		d.SetId("")
	}
	return diag.FromErr(err)
}

func ImportDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	diags := ReadDatabase(ctx, d, meta)

	if diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/saritasa/terraform-provider-mssql/model"
	"github.com/saritasa/terraform-provider-mssql/mssql"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const collationProp = "default_collation"
const nameProp = "name"
const languageProp = "default_language"

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
			nameProp: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			languageProp: {
				Type:     schema.TypeString,
				Optional: true,
				//Default:  "us_english",
			},

			collationProp: {
				Type:     schema.TypeString,
				Optional: true,
				//Default:  "SQL_Latin1_General_CP1_CI_AS",
			},
		},
	}
}

func CreateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)

	database := model.DatabaseFromSchema(d)
	err := connector.CreateDatabase(ctx, database)
	if err == nil {
		d.SetId(database.Name)
	}
	return diag.FromErr(err)
}

func UpdateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	diags := diag.Diagnostics{}

	database := model.DatabaseFromSchema(d)

	if d.HasChange(collationProp) && database.DefaultCollation != "" {
		smtSQL := fmt.Sprintf("ALTER DATABASE %s COLLATE %s", database.Name, database.DefaultCollation)
		err := connector.ExecContext(ctx, smtSQL)
		if err != nil {
			diags = diag.FromErr(err)
		}
	}

	return diags
}

func ReadDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	database := model.DatabaseFromSchema(d)
	err := connector.ReadDatabase(ctx, database)
	if err != nil {
		return diag.FromErr(err)
	}
	return database.ToSchema(d)
}

func DeleteDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	err := connector.DeleteDatabase(ctx, d.Get("name").(string))
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

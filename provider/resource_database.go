package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/saritasa/terraform-provider-mssql/model"
	"github.com/saritasa/terraform-provider-mssql/mssql"

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

			"default_character_set": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
			},

			"default_collation": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8_general_ci",
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
	if err := connector.UpdateDatabase(ctx, model.DatabaseFromSchema(d)); err != nil {
		return diag.FromErr(err)
	}

	return ReadDatabase(ctx, d, meta)
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

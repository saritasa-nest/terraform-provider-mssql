package model

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Database struct {
	Name             string
	DefaultLanguage  string
	DefaultCollation string
}

func DatabaseFromSchema(data *schema.ResourceData) *Database {
	database := &Database{
		Name:             data.Get("Name").(string),
		DefaultLanguage:  data.Get("default_language").(string),
		DefaultCollation: data.Get("default_collation").(string),
	}
	return database
}

func (database *Database) ToSchema(d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := d.Set("Name", database.Name)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("default_language", database.DefaultLanguage)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("default_collation", database.DefaultCollation)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}
	return diags
}

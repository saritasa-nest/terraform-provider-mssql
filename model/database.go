package model

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Database struct {
	Name             string
	DefaultCollation string
	Options          OptionsList
}

func (d *Database) Parse(data *schema.ResourceData) *Database {
	d.Name = data.Get("name").(string)
	d.DefaultCollation = data.Get("default_collation").(string)
	d.Options = make(OptionsList)
	return d
}

func (database *Database) ToSchema(d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := d.Set("Name", database.Name)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("options", database.Options)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("default_collation", database.DefaultCollation)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}
	return diags
}

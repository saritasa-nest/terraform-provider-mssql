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

func (d *Database) ToSchema(data *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := data.Set("name", d.Name)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = data.Set("options", d.Options)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = data.Set("default_collation", d.DefaultCollation)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}
	return diags
}

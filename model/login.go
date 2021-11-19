package model

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Login struct {
	Name     string
	Password string
	Options  OptionsList
}

func (login *Login) Parse(data *schema.ResourceData) *Login {
	login.Name = data.Get("name").(string)
	login.Password = data.Get("password").(string)
	login.Options = make(OptionsList).Parse(data.Get("options").(map[string]interface{}))
	return login
}

func (login *Login) ToSchema(d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := d.Set("name", login.Name)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("options", login.Options)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("password", login.Password)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	return diags
}

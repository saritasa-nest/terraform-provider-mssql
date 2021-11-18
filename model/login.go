package model

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Login struct {
	Name            string
	DefaultDatabase string
	DefaultLanguage string
	Password        string
}

func LoginFromSchema(data *schema.ResourceData) *Login {
	database := &Login{
		Name:            data.Get("name").(string),
		DefaultDatabase: data.Get("default_database").(string),
		DefaultLanguage: data.Get("default_language").(string),
		Password:        data.Get("password").(string),
	}
	return database
}

func (login *Login) ToSchema(d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := d.Set("name", login.Name)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("default_language", login.DefaultLanguage)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("default_database", login.DefaultDatabase)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("password", login.Password)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	return diags
}

package model

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Login struct {
	PrincipalID     int64
	LoginName       string
	DefaultDatabase string
	DefaultLanguage string
}

func LoginFromSchema(data *schema.ResourceData) *Login {
	database := &Login{
		PrincipalID:     int64(data.Get("principal_id").(int)),
		LoginName:       data.Get("login_name").(string),
		DefaultDatabase: data.Get("default_database").(string),
		DefaultLanguage: data.Get("default_language").(string),
	}
	return database
}

func (login *Login) ToSchema(d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := d.Set("login_name", login.LoginName)
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

	err = d.Set("principal_id", int(login.PrincipalID))
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	return diags
}

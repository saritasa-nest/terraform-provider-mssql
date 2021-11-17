package model

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type User struct {
	PrincipalID     int64
	Username        string
	ObjectId        string
	LoginName       string
	Password        string
	AuthType        string
	DefaultSchema   string
	DefaultLanguage string
	Roles           []string
}

func SchemaToUser(d *schema.ResourceData) User {
	user := User{
		PrincipalID:     int64(d.Get("principal_id").(int)),
		Username:        d.Get("username").(string),
		ObjectId:        d.Get("object_id").(string),
		LoginName:       d.Get("login_name").(string),
		Password:        d.Get("password").(string),
		AuthType:        d.Get("auth_type").(string),
		DefaultSchema:   d.Get("default_schema").(string),
		DefaultLanguage: d.Get("default_language").(string),
		Roles:           nil,
	}
	return user
}

func (user *User) ToSchema(d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := d.Set("username", user.Username)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("object_id", user.ObjectId)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("auth_type", user.AuthType)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("default_schema", user.DefaultSchema)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	return diags
}

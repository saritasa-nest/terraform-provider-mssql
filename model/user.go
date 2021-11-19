package model

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type User struct {
	PrincipalID int
	Database    string
	Username    string
	ObjectId    string
	LoginName   string
	Password    string
	AuthType    string
	Options     OptionsList
	Roles       []string
}

func (user *User) Parse(data *schema.ResourceData) *User {
	if user == nil {
		user = &User{}
	}

	user.PrincipalID = data.Get("principal_id").(int)
	user.Database = data.Get("database").(string)
	user.Username = data.Get("username").(string)
	user.ObjectId = data.Get("object_id").(string)
	user.LoginName = data.Get("login_name").(string)
	user.Password = data.Get("password").(string)
	user.AuthType = data.Get("auth_type").(string)
	user.Options = make(OptionsList).Parse(data.Get("options").(map[string]interface{}))
	user.Roles = nil
	return user
}

func (user *User) ToSchema(d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := d.Set("database", user.Database)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	err = d.Set("username", user.Username)
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

	err = d.Set("options", user.Options)
	if err != nil {
		diags = append(diags, diag.FromErr(err)[0])
	}

	return diags
}

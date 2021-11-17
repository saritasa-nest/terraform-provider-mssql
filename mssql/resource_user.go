package mssql

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/saritasa/terraform-provider-mssql/model"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateUser,
		UpdateContext: UpdateUser,
		ReadContext:   ReadUser,
		DeleteContext: DeleteUser,
		Importer: &schema.ResourceImporter{
			StateContext: ImportUser,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "User password",
			},

			"object_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "External object ID",
			},
			"principal_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"login_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User login",
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_language": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func CreateUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*Connector)
	if err := connector.PingContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	user := model.SchemaToUser(d)

	log.Println("Creating user: ", user.Username)

	err := connector.createUser(ctx, connector.Database, &user)

	if err != nil {
		return diag.FromErr(err)
	}

	userId := fmt.Sprintf("%s/%s", connector.Database, user.Username)
	d.SetId(userId)

	return diag.Diagnostics{}
}

func UpdateUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*Connector)
	if err := connector.PingContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	user := model.User{
		PrincipalID:     int64(d.Get("principal_id").(int)),
		Username:        d.Get("username").(string),
		ObjectId:        d.Get("object_id").(string),
		LoginName:       d.Get("login_name").(string),
		Password:        d.Get("password").(string),
		AuthType:        d.Get("auth_type").(string),
		DefaultSchema:   d.Get("default_schema").(string),
		DefaultLanguage: d.Get("default_language").(string),
		Roles:           d.Get("roles").([]string),
	}

	err := connector.updateUser(ctx, connector.Database, &user)
	return diag.FromErr(err)
}

func ReadUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*Connector)

	user, err := connector.getUser(ctx, connector.Database, d.Get("username").(string))
	diags := diag.FromErr(err)

	if user != nil {
		d.SetId(fmt.Sprintf("%s/%s", connector.Database, user.Username))
		diags = user.ToSchema(d)
	}
	return diags
}

func DeleteUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*Connector)
	user := model.SchemaToUser(d)

	err := connector.deleteUser(ctx, connector.Database, user.Username)

	if err == nil {
		d.SetId("")
	}
	return diag.FromErr(err)
}

func ImportUser(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lastSeparatorIndex := strings.LastIndex(d.Id(), "/")

	if lastSeparatorIndex <= 0 {
		return nil, fmt.Errorf("wrong ID format %s (expected database/username)", d.Id())
	}

	username := d.Id()[0:lastSeparatorIndex]
	database := d.Id()[lastSeparatorIndex+1:]

	connector := meta.(*Connector)
	user, err := connector.getUser(ctx, database, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user '%s' not found", d.Id())
	}

	user.ToSchema(d)

	return []*schema.ResourceData{d}, nil
}

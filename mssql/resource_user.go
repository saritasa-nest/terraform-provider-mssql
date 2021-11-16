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
		Read:          ReadUser,
		Delete:        DeleteUser,
		Importer: &schema.ResourceImporter{
			State: ImportUser,
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
		},
	}
}

func CreateUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*Connector)
	if err := connector.PingContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	user := model.User{
		PrincipalID:     d.Get("principal_id").(int64),
		Username:        d.Get("username").(string),
		ObjectId:        d.Get("object_id").(string),
		LoginName:       d.Get("login_name").(string),
		Password:        d.Get("password").(string),
		AuthType:        d.Get("auth_type").(string),
		DefaultSchema:   d.Get("default_schema").(string),
		DefaultLanguage: d.Get("default_language").(string),
		Roles:           nil,
	}

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
		PrincipalID:     d.Get("principal_id").(int64),
		Username:        d.Get("username").(string),
		ObjectId:        d.Get("object_id").(string),
		LoginName:       d.Get("login_name").(string),
		Password:        d.Get("password").(string),
		AuthType:        d.Get("auth_type").(string),
		DefaultSchema:   d.Get("default_schema").(string),
		DefaultLanguage: d.Get("default_language").(string),
		Roles:           nil,
	}

	err := connector.updateUser(ctx, connector.Database, &user)
	return diag.FromErr(err)
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	db, err := meta.(*MsSqlClient).GetDbConn()
	if err != nil {
		return err
	}

	stmtSQL := fmt.Sprintf("SELECT USER FROM mysql.user WHERE USER='%s'",
		d.Get("user").(string))

	log.Println("Executing statement:", stmtSQL)

	rows, err := db.Query(stmtSQL)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() && rows.Err() == nil {
		d.SetId("")
	}
	return rows.Err()
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	db, err := meta.(*MsSqlClient).GetDbConn()
	if err != nil {
		return err
	}

	stmtSQL := fmt.Sprintf("DROP USER '%s'@'%s'",
		d.Get("user").(string),
		d.Get("host").(string))

	log.Println("Executing statement:", stmtSQL)

	_, err = db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

func ImportUser(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lastSeparatorIndex := strings.LastIndex(d.Id(), "@")

	if lastSeparatorIndex <= 0 {
		return nil, fmt.Errorf("wrong ID format %s (expected USER@HOST)", d.Id())
	}

	user := d.Id()[0:lastSeparatorIndex]
	host := d.Id()[lastSeparatorIndex+1:]

	db, err := meta.(*MsSqlClient).GetDbConn()
	if err != nil {
		return nil, err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(1) FROM mysql.user WHERE user = ? AND host = ?", user, host).Scan(&count)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("user '%s' not found", d.Id())
	}

	d.Set("user", user)
	d.Set("host", host)
	d.Set("tls_option", "NONE")

	return []*schema.ResourceData{d}, nil
}

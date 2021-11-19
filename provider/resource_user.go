package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/saritasa/terraform-provider-mssql/model"
	"github.com/saritasa/terraform-provider-mssql/mssql"
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
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "In which database this user will be created",
			},
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
			"login_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "Create user for existing [login] from 'master' database or Windows login name",
				ConflictsWith: []string{"object_id", "principal_id"},
			},

			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "External object ID",
				ConflictsWith: []string{"login_name", "principal_id"},
			},
			"principal_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Description:   "Create user for existing [Principal ID] from 'master' database",
				ConflictsWith: []string{"object_id", "login_name"},
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DATABASE",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					allowedValues := []string{"DATABASE", "INSTANCE", "EXTERNAL"}
					if !model.InArray(val.(string), allowedValues) {
						errs = append(errs, fmt.Errorf("auth_type must be one of: %s", allowedValues))
					}
					return
				},
			},
			"options": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
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
	connector := meta.(*mssql.Connector)
	user := new(model.User).Parse(d)

	userId := fmt.Sprintf("%s/%s", user.Database, user.Username)

	err := connector.CreateUser1(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(userId)

	return nil
}

func UpdateUser(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	if err := connector.PingContext(ctx); err != nil {
		return diag.FromErr(err)
	}
	user := new(model.User).Parse(data)

	err := connector.UpdateUser(ctx, connector.Database, user)
	return diag.FromErr(err)
}

func ReadUser(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)

	user, err := connector.GetUser(ctx, connector.Database, data.Get("username").(string))
	diags := diag.FromErr(err)

	if user != nil {
		data.SetId(fmt.Sprintf("%s/%s", connector.Database, user.Username))
		diags = user.ToSchema(data)
	}
	return diags
}

func DeleteUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	user := new(model.User).Parse(d)

	stmtSQL := fmt.Sprintf("IF EXISTS (SELECT 1 FROM [%s].[sys].[database_principals] WHERE [name] = '%s') "+
		"DROP USER %s", user.Database, user.Username, user.Username)

	log.Printf("Executing statement: %s", stmtSQL)
	err := connector.ExecContext(ctx, stmtSQL)
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

	connector := meta.(*mssql.Connector)
	user, err := connector.GetUser(ctx, database, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user '%s' not found", d.Id())
	}

	user.ToSchema(d)

	return []*schema.ResourceData{d}, nil
}

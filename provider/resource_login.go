package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/saritasa/terraform-provider-mssql/model"
	"github.com/saritasa/terraform-provider-mssql/mssql"
	"strings"
)

func ResourceLogin() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateLogin,
		ReadContext:   ReadLogin,
		DeleteContext: DeleteLogin,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_database": {
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

func CreateLogin(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	login := model.LoginFromSchema(data)
	stmtSQL := "CREATE LOGIN [" + login.Name + "]"
	if login.Password+login.DefaultDatabase+login.DefaultLanguage != "" {
		stmtSQL += " WITH "
		if login.Password != "" {
			stmtSQL += fmt.Sprintf(" PASSWORD = '%s', ", login.Password)
		}
		if login.DefaultDatabase != "" {
			stmtSQL += fmt.Sprintf(" DEFAULT_DATABASE = %s,", login.DefaultDatabase)
		}
		if login.DefaultLanguage != "" {
			stmtSQL += fmt.Sprintf(" DEFAULT_DATABASE = %s,", login.DefaultLanguage)
		}
		stmtSQL = strings.TrimRight(stmtSQL, ",")
	}

	err := connector.ExecContext(ctx, stmtSQL)
	if err == nil {
		data.SetId(login.Name)
	}

	return diag.FromErr(err)
}

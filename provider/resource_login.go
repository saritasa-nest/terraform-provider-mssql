package provider

import (
	"context"
	"database/sql"
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
		UpdateContext: UpdateLogin,
		DeleteContext: DeleteLogin,

		Importer: &schema.ResourceImporter{
			StateContext: ImportLogin,
		},

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

func ReadLogin(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	var login model.Login
	err := connector.QueryRowContext(ctx,
		"SELECT name, default_database_name, default_language_name FROM [master].[sys].[sql_logins] WHERE [name] = @name",
		func(r *sql.Row) error {
			return r.Scan(&login.Name, &login.DefaultDatabase, &login.DefaultLanguage)
		},
		sql.Named("name", data.Get("Name")),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	return login.ToSchema(data)
}

func UpdateLogin(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	login := model.LoginFromSchema(data)
	diags := diag.Diagnostics{}

	if data.HasChange("default_database") {
		stmtSQL := fmt.Sprintf("ALTER LOGIN [%s] WITH DEFAULT_DATABASE %s", login.Name, login.DefaultDatabase)
		err := connector.ExecContext(ctx, stmtSQL)
		if err != nil {
			diags = append(diags, diag.FromErr(err)[0])
		}
	}

	if data.HasChange("default_language") {
		stmtSQL := fmt.Sprintf("ALTER LOGIN [%s] WITH DEFAULT_LANGUAGE %s", login.Name, login.DefaultLanguage)
		err := connector.ExecContext(ctx, stmtSQL)
		if err != nil {
			diags = append(diags, diag.FromErr(err)[0])
		}
	}

	return diags
}

func DeleteLogin(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(*mssql.Connector)
	name := data.Get("Name").(string)

	err := killSessionsForLogin(connector, ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	stmtSQL := fmt.Sprintf("IF EXISTS (SELECT 1 FROM [master].[sys].[sql_logins] WHERE [name] = '%s') DROP LOGIN [%s]", name, name)
	err = connector.ExecContext(ctx, stmtSQL)
	return diag.FromErr(err)
}

func ImportLogin(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	diags := ReadLogin(ctx, data, meta)
	if diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}

	return []*schema.ResourceData{data}, nil
}

func killSessionsForLogin(c *mssql.Connector, ctx context.Context, name string) error {
	cmd := `-- adapted from https://stackoverflow.com/a/5178097/38055
          DECLARE sessionsToKill CURSOR FAST_FORWARD FOR
            SELECT session_id
            FROM sys.dm_exec_sessions
            WHERE login_name = @name
          OPEN sessionsToKill
          DECLARE @sessionId INT
          DECLARE @statement NVARCHAR(200)
          FETCH NEXT FROM sessionsToKill INTO @sessionId
          WHILE @@FETCH_STATUS = 0
          BEGIN
            PRINT 'Killing session ' + CAST(@sessionId AS NVARCHAR(20)) + ' for model ' + @name
            SET @statement = 'KILL ' + CAST(@sessionId AS NVARCHAR(20))
            EXEC sp_executesql @statement
            FETCH NEXT FROM sessionsToKill INTO @sessionId
          END
          CLOSE sessionsToKill
          DEALLOCATE sessionsToKill`
	return c.ExecContext(ctx, cmd, sql.Named("name", name))
}

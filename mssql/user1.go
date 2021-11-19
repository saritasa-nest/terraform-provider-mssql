package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/saritasa/terraform-provider-mssql/model"
	"log"
	"strings"
)

func (c *Connector) CreateUser1(ctx context.Context, user *model.User) error {
	var version string
	err := c.QueryRowContext(ctx, "SELECT @@VERSION", func(row *sql.Row) error { return row.Scan(&version) })
	if err != nil {
		return err
	}

	stmtSQL := fmt.Sprintf("CREATE USER %s ", user.Username)
	if user.AuthType == "DATABASE" && user.LoginName == "" && user.Password == "" {
		return fmt.Errorf("for 'DATABASE' authentication type user password is required")
	}

	if user.LoginName != "" {
		stmtSQL += fmt.Sprintf("FOR LOGIN [%s]", user.LoginName)
	}
	if user.Password != "" {
		stmtSQL += fmt.Sprintf("WITH PASSWORD = '%s'", user.Password)
	}
	if user.AuthType == "EXTERNAL" {
		if strings.Contains(version, "Microsoft SQL Azure") {
			if user.ObjectId != "" {
				stmtSQL += " WITH SID=CONVERT(varchar(64), CAST(CAST(" + user.ObjectId +
					" AS UNIQUEIDENTIFIER) AS VARBINARY(16)), 1), TYPE=E"
			} else {
				stmtSQL += " FROM EXTERNAL PROVIDER"
			}
		}
	}

	if len(user.Options) > 0 {
		if !strings.Contains(stmtSQL, " WITH ") {
			stmtSQL += " WITH "
		}
		for opt := range user.Options {
			value := user.Options[opt].ValueOrSqlNull()
			stmtSQL += fmt.Sprintf(" %s = %s,", opt, value)
		}
		stmtSQL = strings.TrimRight(stmtSQL, ", ")
	}

	log.Printf("Executing statement: %s", stmtSQL)
	err = c.
		setDatabase(user.Database).
		ExecContext(ctx, stmtSQL)
	return err
}

func (c *Connector) GetUserRoles(ctx context.Context, username string) ([]string, error) {
	roles := make([]string, 0)
	err := c.QueryContext(ctx, `select r.name
		from
			master.sys.server_role_members rm
			inner join
			master.sys.server_principals r on r.principal_id = rm.role_principal_id and r.type = 'R'
			inner join
			master.sys.server_principals m on m.principal_id = rm.member_principal_id
		where m.name = @username`,
		func(rows *sql.Rows) error {
			var roleName string
			err := rows.Scan(&roleName)
			if err == nil {
				roles = append(roles, roleName)
			}
			return err
		}, sql.Named("username", username))

	return roles, err
}

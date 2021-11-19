package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/saritasa/terraform-provider-mssql/model"
	"log"
)

func (c *Connector) CreateDatabase(ctx context.Context, d *model.Database) error {
	stmtSQL := databaseConfigSQL("CREATE", d)
	log.Println("Executing statement:", stmtSQL)

	return c.ExecContext(ctx, stmtSQL)
}

func (c *Connector) ReadDatabase(ctx context.Context, d *model.Database) error {
	stmtSQL := "SELECT name, collation_name FROM sys.databases WHERE name LIKE '" + d.Name + "'"

	log.Println("Executing statement:", stmtSQL)
	var collation NullString
	err := c.QueryRowContext(ctx, stmtSQL, func(row *sql.Row) error {
		return row.Scan(&d.Name, &collation)
	})
	if err != nil {
		return fmt.Errorf("read database info: %s", err)
	} else {
		d.DefaultCollation = fmt.Sprintf("%s", collation)
	}

	return nil
}

func (c *Connector) DeleteDatabase(ctx context.Context, name string) error {
	stmtSQL := "DROP DATABASE " + quoteIdentifier(name)
	log.Println("Executing statement:", stmtSQL)
	return c.ExecContext(ctx, stmtSQL)
}

func databaseConfigSQL(verb string, d *model.Database) string {
	var defaultCharsetClause string
	var defaultCollationClause string

	if d.DefaultLanguage != "" {
		defaultCharsetClause = " DEFAULT LANGUAGE = " + quoteIdentifier(d.DefaultLanguage)
	}
	if d.DefaultCollation != "" {
		defaultCollationClause = " COLLATE " + quoteIdentifier(d.DefaultCollation)
	}

	return fmt.Sprintf(
		"%s DATABASE %s %s %s",
		verb,
		quoteIdentifier(d.Name),
		defaultCharsetClause,
		defaultCollationClause,
	)
}

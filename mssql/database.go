package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/saritasa/terraform-provider-mssql/model"
	"log"
	"strings"
)

const defaultCharacterSetKeyword = "CHARACTER SET "
const defaultCollateKeyword = "COLLATE "

func (c *Connector) CreateDatabase(ctx context.Context, d *model.Database) error {
	stmtSQL := databaseConfigSQL("CREATE", d)
	log.Println("Executing statement:", stmtSQL)

	db, err := c.db()
	if err != nil {
		return err
	}

	_, err = db.Exec(stmtSQL)
	return err
}

func (c *Connector) ReadDatabase(ctx context.Context, d *model.Database) error {
	db, err := c.db()
	if err != nil {
		return err
	}

	stmtSQL := "SELECT name, collation_name FROM sys.databases WHERE name LIKE '" + d.Name + "'"

	log.Println("Executing statement:", stmtSQL)
	var createSQL, _database string
	err = db.QueryRow(stmtSQL).Scan(&_database, &createSQL)
	if err != nil {
		return fmt.Errorf("read database info: %s", err)
	}

	return c.QueryRowContext(ctx, stmtSQL, func(row *sql.Row) error {
		return row.Scan(&d.Name, &d.DefaultCollation)
	})
}

func (c *Connector) UpdateDatabase(ctx context.Context, d *model.Database) error {
	db, err := c.db()
	if err != nil {
		return err
	}

	stmtSQL := databaseConfigSQL("ALTER", d)
	log.Println("Executing statement:", stmtSQL)

	_, err = db.Exec(stmtSQL)
	return err
}

func (c *Connector) DeleteDatabase(ctx context.Context, name string) error {
	db, err := c.db()
	if err != nil {
		return err
	}

	stmtSQL := "DROP DATABASE " + quoteIdentifier(name)
	log.Println("Executing statement:", stmtSQL)

	_, err = db.Exec(stmtSQL)
	return err
}

func databaseConfigSQL(verb string, d *model.Database) string {
	var defaultCharsetClause string
	var defaultCollationClause string

	if d.DefaultLanguage != "" {
		defaultCharsetClause = defaultCharacterSetKeyword + quoteIdentifier(d.DefaultLanguage)
	}
	if d.DefaultCollation != "" {
		defaultCollationClause = defaultCollateKeyword + quoteIdentifier(d.DefaultCollation)
	}

	return fmt.Sprintf(
		"%s DATABASE %s %s %s",
		verb,
		quoteIdentifier(d.Name),
		defaultCharsetClause,
		defaultCollationClause,
	)
}

func extractIdentAfter(sql string, keyword string) string {
	charsetIndex := strings.Index(sql, keyword)
	if charsetIndex != -1 {
		charsetIndex += len(keyword)
		remain := sql[charsetIndex:]
		spaceIndex := strings.IndexRune(remain, ' ')
		return remain[:spaceIndex]
	}

	return ""
}

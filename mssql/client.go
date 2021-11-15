package mssql

import (
	"database/sql"
	"fmt"
	"github.com/hashicorp/go-version"
	"strings"
)

type MsSqlClient struct {
	Host     string
	Port     int
	Username string
	Password string
	db       *sql.DB
}

func GetDbConn(c *MsSqlClient) (*sql.DB, error) {
	if c.db == nil {
		db, err := ConnectToMySQL(c)
		if err != nil {
			return nil, err
		}
		c.db = db
	}
	return c.db, nil
}

func ConnectToMySQL(conf *MsSqlClient) (*sql.DB, error) {
	return nil, nil
}

var identQuoteReplacer = strings.NewReplacer("`", "``")

func QuoteIdentifier(in string) string {
	return fmt.Sprintf("`%s`", identQuoteReplacer.Replace(in))
}

func ServerVersion(db *sql.DB) (*version.Version, error) {
	var versionString string
	err := db.QueryRow("SELECT @@GLOBAL.innodb_version").Scan(&versionString)
	if err != nil {
		return nil, err
	}

	return version.NewVersion(versionString)
}

func ServerVersionString(db *sql.DB) (string, error) {
	var versionString string
	err := db.QueryRow("SELECT @@GLOBAL.version").Scan(&versionString)
	if err != nil {
		return "", err
	}

	return versionString, nil
}

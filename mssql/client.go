package mssql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type MsSqlClient struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Timeout  time.Duration
	db       *sql.DB
}

func (client *MsSqlClient) GetDbConn() (*sql.DB, error) {
	if client.db == nil {
		db, err := ConnectToMySQL(client)
		if err != nil {
			return nil, err
		}
		client.db = db
	}
	return client.db, nil
}

func ConnectToMySQL(conf *MsSqlClient) (*sql.DB, error) {
	return nil, nil
}

var identQuoteReplacer = strings.NewReplacer("`", "``")

func QuoteIdentifier(in string) string {
	return fmt.Sprintf("`%s`", identQuoteReplacer.Replace(in))
}

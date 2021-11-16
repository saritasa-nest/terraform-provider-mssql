package mssql

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const defaultCharacterSetKeyword = "CHARACTER SET "
const defaultCollateKeyword = "COLLATE "
const unknownDatabaseErrCode = 1049

func ResourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateDatabase,
		Update:        UpdateDatabase,
		Read:          ReadDatabase,
		Delete:        DeleteDatabase,
		Importer: &schema.ResourceImporter{
			State: ImportDatabase,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"default_character_set": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
			},

			"default_collation": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8_general_ci",
			},
		},
	}
}

func CreateDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db, err := meta.(*Connector).db()
	if err != nil {
		return diag.FromErr(err)
	}

	stmtSQL := databaseConfigSQL("CREATE", d)
	log.Println("Executing statement:", stmtSQL)

	_, err = db.Exec(stmtSQL)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("name").(string))

	return diag.Diagnostics{}
}

func UpdateDatabase(d *schema.ResourceData, meta interface{}) error {
	db, err := meta.(*Connector).db()
	if err != nil {
		return err
	}

	stmtSQL := databaseConfigSQL("ALTER", d)
	log.Println("Executing statement:", stmtSQL)

	_, err = db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	return ReadDatabase(d, meta)
}

func ReadDatabase(d *schema.ResourceData, meta interface{}) error {
	db, err := meta.(*Connector).db()
	if err != nil {
		return err
	}

	// This is kinda flimsy-feeling, since it depends on the formatting
	// of the SHOW CREATE DATABASE output... but this data doesn't seem
	// to be available any other way, so hopefully MySQL keeps this
	// compatible in future releases.

	name := d.Id()
	stmtSQL := "SHOW CREATE DATABASE " + quoteIdentifier(name)

	log.Println("Executing statement:", stmtSQL)
	var createSQL, _database string
	err = db.QueryRow(stmtSQL).Scan(&_database, &createSQL)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == unknownDatabaseErrCode {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error during show create database: %s", err)
	}

	defaultCharset := extractIdentAfter(createSQL, defaultCharacterSetKeyword)
	defaultCollation := extractIdentAfter(createSQL, defaultCollateKeyword)

	// TODO

	d.Set("name", name)
	d.Set("default_character_set", defaultCharset)
	d.Set("default_collation", defaultCollation)

	return nil
}

func DeleteDatabase(d *schema.ResourceData, meta interface{}) error {
	db, err := meta.(*Connector).db()
	if err != nil {
		return err
	}

	name := d.Id()
	stmtSQL := "DROP DATABASE " + quoteIdentifier(name)
	log.Println("Executing statement:", stmtSQL)

	_, err = db.Exec(stmtSQL)
	if err == nil {
		d.SetId("")
	}
	return err
}

func databaseConfigSQL(verb string, d *schema.ResourceData) string {
	name := d.Get("name").(string)
	defaultCharset := d.Get("default_character_set").(string)
	defaultCollation := d.Get("default_collation").(string)

	var defaultCharsetClause string
	var defaultCollationClause string

	if defaultCharset != "" {
		defaultCharsetClause = defaultCharacterSetKeyword + quoteIdentifier(defaultCharset)
	}
	if defaultCollation != "" {
		defaultCollationClause = defaultCollateKeyword + quoteIdentifier(defaultCollation)
	}

	return fmt.Sprintf(
		"%s DATABASE %s %s %s",
		verb,
		quoteIdentifier(name),
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

func ImportDatabase(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := ReadDatabase(d, meta)

	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

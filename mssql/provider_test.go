package mssql

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// To run these acceptance tests, you will need access to a MS SQL server.
// Amazon RDS is one way to get a MS SQL server. If you use RDS, you can
// use the root account credentials you specified when creating an RDS
// instance to get the access necessary to run these tests. (the tests
// assume full access to the server.)
//
// Set the MSSQL_ENDPOINT and MSSQL_USERNAME environment variables before
// running the tests. If the given user has a password then you will also need
// to set MSSQL_PASSWORD.
//
// The tests assume a reasonably-vanilla MS SQL configuration. In particular,
// they assume that the "utf8" character set is available and that
// "utf8_bin" is a valid collation that isn't the default for that character
// set.
//
// You can run the tests like this:
//    make testacc TEST=./builtin/providers/mysql

var TestProviderFactories map[string]func() (*schema.Provider, error)
var TestAccProvider *schema.Provider

func init() {
	TestAccProvider = Provider()
	TestProviderFactories = map[string]func() (*schema.Provider, error){
		"mssql": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func testProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"MSSQL_ENDPOINT", "MSSQL_USERNAME"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("MSSQL_ENDPOINT, MSSQL_USERNAME and optionally MSSQL_PASSWORD must be set for acceptance tests")
		}
	}

	err := TestAccProvider.Configure(context.TODO(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}

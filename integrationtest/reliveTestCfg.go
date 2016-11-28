package integrationtest

import (
	"database/sql"
	"fmt"

	"github.com/msproject/relive/testtools"
)

// TestCfg - config params used for integration test
type TestCfg struct {
	reliveServerURL string
	mysqlAccessAddr string
}

var (
	reliveTestCfg = &TestCfg{
		reliveServerURL: "http://localhost:9999",
	}
)

func reliveSetupTestEnv() error {
	return nil
}

func cleanupAllTables(dbconnect string) error {
	db, err := sql.Open("mysql", dbconnect)
	if err != nil {
		return fmt.Errorf("could not connect to mysql for cleanup %s", err.Error())
	}

	err = testtools.CleanUpTables(db)
	if err != nil {
		return fmt.Errorf("could not cleanup mysql tables: %s", err.Error())
	}

	return nil
}

package dbinit

import (
	"../logger"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
)

// Config - params needed for DB init to run
type Config struct {
	metaDataURL string
	createDDLs  []string
	deleteDDLs  []string
	logObj      *logger.Logger
}

// NewDBInitConfig - create a new Config object
func NewDBInitConfig(metaURL string, createDDLStmts, deleteDDLStmts []string, logObject *logger.Logger) (*Config, error) {

	if len(metaURL) < 1 {
		return &Config{}, fmt.Errorf("Error: expected env variable 'DSN', found none")
	}

	dbmcfg := &Config{
		metaDataURL: metaURL,
		createDDLs:  createDDLStmts,
		deleteDDLs:  deleteDDLStmts,
		logObj:      logObject,
	}
	return dbmcfg, nil
}

// Configure - configure (or run) the DB schema statements
// (1) CreatePhase  - run CREATE and ALTER DDLs
// (2) DeletePhase - run DELETE/DROP DDLs
func (d *Config) Configure() error {
	d.logObj.PrintInfo("Config.Configure()")

	err := d.createPhase()
	if err != nil {
		return err
	}

	err = d.deletePhase()
	if err != nil {
		return err
	}

	return nil
}

func (d *Config) createPhase() error {
	err := d.createDataBase()
	if err != nil {
		return fmt.Errorf("Create DB failed, err[%s]", err.Error())
	}

	/* A new db connection is used after createDatabase */
	db, err := sql.Open("mysql", d.metaDataURL)
	if err != nil {
		return fmt.Errorf("could not open DB %s err[%s]", d.metaDataURL, err.Error())
	}

	errors := d.getAndExecStmts(db)
	if errors != nil {
		return errors
	}
	db.Close()
	return nil
}

func (d *Config) deletePhase() error {

	if len(d.deleteDDLs) <= 0 {
		d.logObj.PrintInfo("No statements to execute in Delete DDLs")
		return nil
	}
	/* A new db connection is used for Delete Stmts */
	db, err := sql.Open("mysql", d.metaDataURL)
	if err != nil {
		return fmt.Errorf("could not open DB %s err[%s]", d.metaDataURL, err.Error())
	}

	errors := d.execSQLStmts(db, d.deleteDDLs)
	if errors != nil {
		return errors
	}
	db.Close()
	return nil
}

func (d *Config) createDataBase() error {
	var baseMetaURL string
	/* first, strip off the dbName and options part (if it exists) from metaURL
	 * to get the base URL. use this to create database, and then reconstruct the
	 * metaURL with appropriate dbName and options. Note that we reconstruct
	 * metadataURL only if incoming metaURL does not have the dbName
	 * and/or options */
	baseMetaURL = d.metaDataURL[:(strings.Index(d.metaDataURL, "/") + 1)]

	if len(d.metaDataURL) <= (strings.Index(d.metaDataURL, "/") + 1) {
		/* here there is no DBname and/or options specified in input
		 * reconstruct d.metaDataURL to have:
		 * - dbname from USE DDL and
		 * - default options: interpolateParams=true, parseTime=true
		 * Note that we MUST mandate to have a USE DDL in the createDDLs */
		var dbNameWithOptions string
		for _, stmt := range d.createDDLs {
			if strings.Contains(stmt, "USE") {
				dbNameWithOptions = stmt[(len("USE") + 1):(len(stmt) - 1)]
				break
			}
		}

		if len(dbNameWithOptions) <= 0 {
			return fmt.Errorf("DB name is empty. USE stmt does not exist. Cannot proceed to further configuration")
		}

		dbNameWithOptions = dbNameWithOptions + "?interpolateParams=true&parseTime=true"
		d.metaDataURL = d.metaDataURL + dbNameWithOptions
	}

	for _, stmt := range d.createDDLs {
		/* only run "create database" DDL and break out of the loop */
		if strings.Contains(stmt, "CREATE DATABASE") {
			db, err := sql.Open("mysql", baseMetaURL)
			if err != nil {
				return fmt.Errorf("could not connect to sql DB using %s %s", baseMetaURL, err.Error())
			}
			_, err = db.Exec(stmt)
			if err != nil {
				fmt.Println(err.Error())
				if driverErr, ok := err.(*mysql.MySQLError); ok {
					/* ignore SQL errors:
					 * 1. Error 1007: Creating an existing database
					 * 2. Error 1008: Dropping a non-existent database
					 */
					if driverErr.Number == 1007 || driverErr.Number == 1008 {
						err = nil
					} else {
						err = fmt.Errorf("error creating database %s", err.Error())
					}
				} else {
					err = fmt.Errorf("error creating database %s", err.Error())
				}
				if err != nil {
					return err
				}
			}
			db.Close()
			break
		}
	}
	return nil
}

func (d *Config) getAndExecStmts(db *sql.DB) (errors error) {
	/* Order of executing DDLs is important because
	 * unless we CREATE TABLEs we cannot CREATE VIEWs.
	 * ALTER comes afterwards logically, because first
	 * version of a DDL schema cannot contain any ALTER
	 * stmts. Hence the stmtTypes is constructed to run
	 * the DDL statements in this particular order. */
	stmtTypes := []string{"CREATE TABLE",
		"CREATE VIEW",
		"ALTER",
	}
	for _, stmtType := range stmtTypes {
		stmtSet := d.getSQLStmts(stmtType)
		if len(stmtSet) > 0 {
			errors := d.execSQLStmts(db, stmtSet)
			if errors != nil {
				return errors
			}
		}
	}
	return nil
}

func (d *Config) execSQLStmts(db *sql.DB, stmtSet []string) (err error) {
	for _, stmt := range stmtSet {
		_, err := db.Exec(stmt)
		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				/* ignore SQL errors:
				 * 1. Error 1050: Adding existing tables
				 * 2. Error 1051: Dropping non-existent table
				 * 3. Error 1060: Adding existing columns
				 * 4. Error 1091: Dropping a non-existent column
				 */
				if driverErr.Number == 1050 || driverErr.Number == 1051 ||
					driverErr.Number == 1060 || driverErr.Number == 1091 {
					d.logObj.PrintInfo("Ignore a DB error, considered normal during DB init: [%s]", err.Error())
					err = nil
				} else {
					fmt.Println(err.Error())
					err = fmt.Errorf("error issuing sql statement %s %s", stmt, err.Error())
					break
				}
			} else {
				fmt.Println(err.Error())
				err = fmt.Errorf("error issuing sql statement %s %s", stmt, err.Error())
				break
			}
		}
	}

	return err
}

func (d *Config) getSQLStmts(stmtType string) (result []string) {
	for _, stmt := range d.createDDLs {
		if strings.Contains(stmt, stmtType) {
			result = append(result, stmt)
		}
	}
	return result
}

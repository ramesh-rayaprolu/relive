package main

import (
	"./dbi"
	"./dbinit"
	"./dbmodel"
	"./logger"
	"flag"
	"os"
	"time"
)

func main() {
	var metaURL string
	var dbTimeout time.Duration
	flag.StringVar(&metaURL, "metaurl", "root:@tcp(127.0.0.1:9306)/relive?parseTime=true&interpolateParams=true", "URL of the metadata service")
	flag.DurationVar(&dbTimeout, "dbtimeout", 10*time.Second, "timeout for DB queries")
	flag.Parse()

	logObj, _ := logger.NewLoggerObject(false)
	/* first DBInit */
	var dbInitCfg *dbinit.Config
	var err error
	if dbInitCfg, err = dbinit.NewDBInitConfig(metaURL, dbmodel.TableCreateSQL, dbmodel.TableDeleteSQL, logObj); err != nil {
		logObj.PrintError("DB Init failed, exiting. Error: %v", err)
		os.Exit(-1)
	}

	if err = dbInitCfg.Configure(); err != nil {
		logObj.PrintError("DB Init failed, exiting. Error: %v", err)
		os.Exit(-1)
	}
	/* end of DBInit */

	sqlDbi, err := dbi.InitializeDBI(metaURL, dbTimeout)

	if err != nil {
		logObj.PrintError("Could not initialize memsql Dbi %s error %s", metaURL, err.Error())
		os.Exit(1)
	}

	accountLoader := sqlDbi.AccountDBI

	acEntry := &dbmodel.AccountEntry{
		ID:           1,
		PID:          11,
		FirstName:    "Pratibha",
		LastName:     "Revankar",
		EmailID:      "pratirvce@gmail.com",
		PasswdDigest: "pratirvce",
		Role:         0,
	}

	err = accountLoader.AddAccounts(acEntry)

	if err != nil {
		logObj.PrintError("error executing DB: %v", err)
	}
}

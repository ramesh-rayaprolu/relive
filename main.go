package main

import (
	"flag"
	"fmt"
	"github.com/msproject/relive/api"
	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/dbinit"
	"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/logger"
	"net/http"
	"os"
	"time"
)

func main() {
	var metaURL, listen, listenSSL, certFilePath, keyFilePath string
	var dbTimeout time.Duration
	flag.StringVar(&metaURL, "metaurl", "root:@tcp(127.0.0.1:3306)/relive?parseTime=true&interpolateParams=true", "URL of the metadata service")
	flag.DurationVar(&dbTimeout, "dbtimeout", 10*time.Second, "timeout for DB queries")
	flag.StringVar(&listen, "listen", ":9999", "Host and HTTP port to listen on")
	flag.StringVar(&listenSSL, "listenssl", ":8443", "Host and HTTPS port to listen on")
	flag.StringVar(&certFilePath, "cert", "./relive_cert.pem", "absolute file path for the SSL certificate file")
	flag.StringVar(&keyFilePath, "key", "./relive_key.pem", "absolute file path for the SSL key file")
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

	sqlDbi, err := dbi.InitializeDBI(metaURL, dbTimeout, logObj)

	if err != nil {
		logObj.PrintError("Could not initialize the SQL Dbi %s error %s", metaURL, err.Error())
		os.Exit(1)
	}

	//if the account table exists, make sure there exists a root user, if not create it
	initErr := api.InitAccountsDB(sqlDbi)
	if initErr != nil {
		logObj.PrintError("Could not initialize Accounts table. Error %s", initErr.Error())
		os.Exit(1)
	}

	productInitErr := api.InitProductsDB(sqlDbi)
	if productInitErr != nil {
		logObj.PrintError("Could not initialize Products table. Error %s", productInitErr.Error())
		os.Exit(1)
	}

	accountAPI := api.AccountsAPI{
		AccountDBI:      sqlDbi.AccountDBI,
		SubscriptionDBI: sqlDbi.SubscriptionDBI,
		LogObj:          logObj,
	}

	/*productAPI := api.ProductsAPI{
		ProductDBI: sqlDbi.ProductDBI,
		//SubscriptionDBI: sqlDbi.SubscriptionDBI,
		LogObj: logObj,
	}*/

	subscriptionAPI := api.SubscriptionAPI{
		SubscriptionDBI:        sqlDbi.SubscriptionDBI,
		SubscriptionAccountDBI: sqlDbi.SubscriptionAccountDBI,
		ProductDBI:             sqlDbi.ProductDBI,
		LogObj:                 logObj,
	}
	paymentAPI := api.PaymentAPI{
		PaymentDBI:        sqlDbi.PaymentDBI,
		PaymentHistoryDBI: sqlDbi.PaymentHistoryDBI,
		LogObj:            logObj,
	}
	mediaAPI := api.MediaAPI{
		MediaDBI: sqlDbi.MediaTypeDBI,
		LogObj:   logObj,
	}

	router := api.Router{
		Account:      accountAPI,
		Subscription: subscriptionAPI,
		Payment:      paymentAPI,
		Media:        mediaAPI,
		AccountDBI:   sqlDbi.AccountDBI,
		LogObj:       logObj,
	}

	routerSSL := api.RouterSSL{
		Account:      accountAPI,
		Subscription: subscriptionAPI,
		Payment:      paymentAPI,
		Media:        mediaAPI,
		LogObj:       logObj,
	}

	logObj.PrintInfo("Listening on (HTTP) : %s\n", listen)
	logObj.PrintInfo("Listening on (HTTPS): %s\n", listenSSL)

	/* HTTP Server MUX */
	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", router)
	httpMux.HandleFunc("/version", VersionHandler)
	httpMux.HandleFunc("/health", HealthHandler)
	//  Start HTTP
	go func() {
		err := http.ListenAndServe(listen, httpMux)
		if err != nil {
			logObj.PrintError("error Listening : %v", err)
			os.Exit(1)
		}
	}()

	/* HTTPS Server MUX */
	httpsMux := http.NewServeMux()
	httpsMux.Handle("/api/", routerSSL)
	httpsMux.HandleFunc("/version", VersionHandler)
	httpsMux.HandleFunc("/health", HealthHandler)
	//  Start HTTPS
	err = http.ListenAndServeTLS(listenSSL, certFilePath, keyFilePath, httpsMux)
	if err != nil {
		logObj.PrintError("error Listening : %v", err)
		os.Exit(1)
	}
}

// VersionHandler will print the generated build info when /version is requested.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", "1.0.0")
	return
}

// HealthHandler simply returns 204 not content
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Authorization,DNT,X-Auth,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		w.Header().Set("Content-Type", "text/plain charset=UTF-8")
		w.Header().Set("Content-Length", "0")
		return
	}
	if r.Method == "POST" || r.Method == "GET" || r.Method == "PUT" || r.Method == "DELETE" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Authorization,DNT,X-Auth,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	}
	w.WriteHeader(http.StatusNoContent)
}

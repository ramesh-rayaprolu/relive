package api

import (
	"fmt"
	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/logger"
	"net/http"
	"strings"
)

//Router - main HTTP handler for all relive non-SSL APIs
type Router struct {
	Account      http.Handler
	Subscription http.Handler
	Payment      http.Handler
	Media        http.Handler
	Product      http.Handler
	AccountDBI   dbi.AccountTblDBI
	LogObj       *logger.Logger
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Authorization,DNT,X-Auth,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		w.Header().Set("Content-Type", "text/plain charset=UTF-8")
		w.Header().Set("Content-Length", "0")
		return
	}
	if req.Method == "POST" || req.Method == "GET" || req.Method == "PUT" || req.Method == "DELETE" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Authorization,DNT,X-Auth,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	}

	url := req.URL.String()

	r.LogObj.PrintInfo("request URL: %s", url)

	if strings.Contains(url, "/accounts/") || strings.Contains(url, "/accounts?") {
		r.Account.ServeHTTP(w, req)
		return
	} else if strings.Contains(url, "/subscription/") || strings.Contains(url, "/subscription?") {
		r.Subscription.ServeHTTP(w, req)
		return
	} else if strings.Contains(url, "/payment/") || strings.Contains(url, "/payment?") {
		r.Payment.ServeHTTP(w, req)
		return
	} else if strings.Contains(url, "/media/") || strings.Contains(url, "/media?") {
		r.Media.ServeHTTP(w, req)
		return
	} else if strings.Contains(url, "/products/") || strings.Contains(url, "/products?") {
		r.Product.ServeHTTP(w, req)
		return
	}

	http.Error(w, fmt.Sprintf("Invalid Request made. %s is not an active endpoint.", url), http.StatusMethodNotAllowed)
}

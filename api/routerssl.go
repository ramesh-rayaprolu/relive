package api

import (
	"../dbi"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

//RouterSSL - main HTTP handler for all relive SSL APIs
type RouterSSL struct {
	Account      http.Handler
	Subscription http.Handler
	Payment      http.Handler
	Media        http.Handler
	AccountsDBI  dbi.AccountTblDBI
}

func (r RouterSSL) ServeHTTP(w http.ResponseWriter, req *http.Request) {

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

	if strings.Contains(url, "/accounts/login") || strings.Contains(url, "/accounts/register") {
		r.Account.ServeHTTP(w, req)
		return
	}

	/* Authenticate */
	_, err := authenticate(r.AccountsDBI, w, req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	/*if role < 1 && (strings.Contains(url, "/accounts/search") || strings.Contains(url, "/recordings/update") ||
		strings.Contains(url, "/recordings/delete") || strings.Contains(url, "/streams/configure") ||
		strings.Contains(url, "/streams/create") || strings.Contains(url, "/streams/update") ||
		strings.Contains(url, "/streams/streamstate") || strings.Contains(url, "/streams/delete")) {
		http.Error(w, "Insufficient role privilege!", http.StatusForbidden)
		return
	} else if role < 2 && (strings.Contains(url, "/accounts/create") || strings.Contains(url, "/accounts/update") ||
		strings.Contains(url, "/accounts/delete") || strings.Contains(url, "/accounts/activate") ||
		strings.Contains(url, "/accounts/disable")) {
		http.Error(w, "Insufficient role privilege!", http.StatusForbidden)
		return
	} else */
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
	}

	http.Error(w, fmt.Sprintf("Invalid Request made. %s is not an active endpoint.", url), http.StatusMethodNotAllowed)
}

func authenticate(accountDBI dbi.AccountTblDBI, w http.ResponseWriter, r *http.Request) (int, error) {
	authData := r.Header.Get("Authorization")
	if authData == "" {
		return 0, fmt.Errorf("authorization header is empty")
	}
	// Decode authData
	authEncoded := strings.TrimPrefix(authData, "Basic ")
	authDecoded, _ := base64.URLEncoding.DecodeString(authEncoded)
	authArray := strings.Split(string(authDecoded), ":")

	loginResult, err := accountDBI.Login(authArray[0], authArray[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, err
	}
	if loginResult == nil {
		http.Error(w, "Account does not exist.", http.StatusUnauthorized)
		return 0, fmt.Errorf("account does not exist")
	}
	role := loginResult.Role
	return role, nil
}

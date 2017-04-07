package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/logger"
	"github.com/msproject/relive/util"
)

// AccountsAPI struct
type AccountsAPI struct {
	AccountDBI      dbi.AccountTblDBI
	SubscriptionDBI dbi.SubscriptionTblDBI
	LogObj          *logger.Logger
}

//InitAccountsDB - create root user if it doesnt exist
func InitAccountsDB(sqlDBI dbi.DBI) (err error) {
	var exists bool

	exists, err = sqlDBI.AccountDBI.CheckAccountTableExists()

	if err != nil {
		fmt.Printf("Error checking for Account Table < %s >\n", err.Error())
		return err
	}

	if !exists {
		fmt.Println("AccountTable Does not exist")
		return nil
	}

	exists, err = sqlDBI.AccountDBI.CheckAccountExists("root")
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	req := util.CreateAccountReq{
		UserName:  "root",
		Email:     "root@relive.com",
		FirstName: "root",
		LastName:  "root",
		PWD:       "video@Cloud",
		Role:      0,
		CompanyID: 0,
	}
	err = sqlDBI.AccountDBI.CreateAccount(req)
	if err != nil {
		fmt.Printf("Error creating root account %s\n", err.Error())
		return err
	}
	return nil
}

//	/api/accounts/search  updated
//func handleAccountsSearch(UserName string) error {
func handleAccountsSearch(api AccountsAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	var username string
	var req util.SearchAccountReq
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/accounts/Search")
	}

	URLSuffix := args[0]
	parsedURLSuffix, err := url.Parse(URLSuffix)
	params := parsedURLSuffix.Query()

	if len(params["user"]) > 0 {
		username = params["user"][0]
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required query parameters NOT specified in search request")
	}

	req, err = api.AccountDBI.SearchAccount(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	err = writeResponse(req, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

// /api/accounts/create
func handleAccountsCreate(api AccountsAPI, args []string, w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/accounts/create")
	}

	var req util.CreateAccountReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	if req.Email == "" || req.UserName == "" || req.LastName == "" || req.FirstName == "" || req.PWD == "" {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required parameters NOT specified in search request")
	}

	exists, err1 := api.AccountDBI.CheckAccountExists(req.UserName)
	if err1 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err1
	}
	if exists {
		http.Error(w, "Account exists.", http.StatusUnauthorized)
		return err1
	}

	err := api.AccountDBI.CreateAccount(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	//fmt.Println("Inside search")
	return nil
}

// /api/accounts/update - Update an account
func handleAccountsUpdate(api AccountsAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/accounts/my/update - Update my account
func handleMyAccountUpdate(api AccountsAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/accounts/change - change password
func handleChangePassword(api AccountsAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/accounts/login/  - login
func handleAccountsLogin(api AccountsAPI, args []string, w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/accounts/login")
	}

	var req util.LoginReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	exists, err1 := api.AccountDBI.CheckAccountExists(req.UserName)
	if err1 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err1
	}
	if !exists {
		http.Error(w, "Account does not exist.", http.StatusUnauthorized)
		return err1
	}

	recs, err := api.AccountDBI.Login(req.UserName, req.PWD)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	err = writeResponse(recs, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func (api AccountsAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, d := range account {
		if d.re.MatchString(r.URL.String()) {
			err := d.f(api, d.re.FindStringSubmatch(r.URL.String()), w, r)
			if err != nil {
				returnMessage := fmt.Sprintf("%v", err)
				w.Write([]byte(returnMessage))
			}
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("No match found.\n"))
}

type accountT struct {
	regex string
	re    *regexp.Regexp
	f     func(api AccountsAPI, args []string, w http.ResponseWriter, r *http.Request) error
}

var account []accountT

func init() {
	var regex string
	regex = "/api/accounts/search\\?([^/]+)$"
	account = append(account,
		accountT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleAccountsSearch,
		},
	)
	regex = "/api/accounts/create$"
	account = append(account,
		accountT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleAccountsCreate,
		},
	)
	regex = "/api/accounts/update$"
	account = append(account,
		accountT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleAccountsUpdate,
		},
	)
	regex = "/api/accounts/my/update$"
	account = append(account,
		accountT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMyAccountUpdate,
		},
	)
	regex = "/api/accounts/change$"
	account = append(account,
		accountT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleChangePassword,
		},
	)
	regex = "/api/accounts/login$"
	account = append(account,
		accountT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleAccountsLogin,
		},
	)
}

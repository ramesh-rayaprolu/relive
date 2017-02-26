package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/logger"
	"github.com/msproject/relive/util"
)

// SubscriptionAPI struct
type SubscriptionAPI struct {
	SubscriptionDBI        dbi.SubscriptionTblDBI
	SubscriptionAccountDBI dbi.SubscriptionAccountTblDBI
	ProductDBI             dbi.ProductTblDBI
	LogObj                 *logger.Logger
}

//	/api/subscription/search
func handleSubscriptionSearch(api SubscriptionAPI, args []string, w http.ResponseWriter, r *http.Request) error {

	// check for API Method
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/subscription/search")
	}

	return nil
}

// /api/subscription/create
func handleSubscriptionCreate(api SubscriptionAPI, args []string, w http.ResponseWriter, r *http.Request) error {

	// check for API Method
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/subscription/create")
	}

	// decode the JSON against the structure
	var req util.CreateSubscriptionReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	err := api.SubscriptionDBI.CreateSubscription(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

// /api/subscription/update -
func handleSubscriptionUpdate(api SubscriptionAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	// check for API Method
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/subscription/update")
	}

	// decode the JSON against the structure
	var req util.CreateSubscriptionReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	fmt.Println("\nSubscriptionCode and NumberOfAdmins: ", req.SubscriptionCode, req.NumberOfAdmins)

	if req.SubscriptionCode == 0 || req.NumberOfAdmins == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required parameters NOT specified in update request")
	}

	fmt.Println("\nAfter the check")

	err := api.SubscriptionDBI.UpdateSubscription(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

// /api/subscription/delete -
func handleSubscriptionDelete(api SubscriptionAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	// check for API Method
	if r.Method != "DELETE" {
		w.Header().Set("Allow", "DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/subscription/delete")
	}

	return nil
}

type subscriptionT struct {
	regex string
	re    *regexp.Regexp
	f     func(api SubscriptionAPI, args []string, w http.ResponseWriter, r *http.Request) error
}

var subscription []subscriptionT

func (api SubscriptionAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, d := range subscription {
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

func init() {
	var regex string
	regex = "/api/subscription/search$"
	subscription = append(subscription,
		subscriptionT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleSubscriptionSearch,
		},
	)
	regex = "/api/subscription/create$"
	subscription = append(subscription,
		subscriptionT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleSubscriptionCreate,
		},
	)
	regex = "/api/subscription/update$"
	subscription = append(subscription,
		subscriptionT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleSubscriptionUpdate,
		},
	)
	regex = "/api/subscription/delete$"
	subscription = append(subscription,
		subscriptionT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleSubscriptionDelete,
		},
	)
}

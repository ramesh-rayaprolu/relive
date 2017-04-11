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

	fmt.Println("before prtining args")
	/* anything after /api/subscription/search/ will be in args, split by '/' */
	fmt.Println("arguments in search: ", args)

	//err := api.SubscriptionDBI.SearchSubscription(strconv.Atoi(args))
	var subscrCode uint32
	subscrCode = 3
	fmt.Println("subscrCode = ", args)

	var subs []util.SubscrDetails
	subs, err := api.SubscriptionDBI.SearchSubscription(subscrCode)
	//err := api.SubscriptionDBI.SearchSubscription(uint32(args))
	//fmt.Println(subs)

	jsonStr, err := json.Marshal(subs)
	fmt.Println("json: ", jsonStr)

	w.Header().Set("Content-Type", "application/json")
	n, err := w.Write(jsonStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Failure to write, err = %s", err)
	}
	if n != len(jsonStr) {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Short write sent = %d, wrote = %d", len(jsonStr), n)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
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

	if req.SubscriptionCode == 0 || req.NumberOfAdmins == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required parameters NOT specified in update request")
	}

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

	// decode the JSON against the structure
	var req util.CreateSubscriptionReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	if req.SubscriptionCode == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required parameters NOT specified in delete request")
	}

	err := api.SubscriptionDBI.DeleteSubscription(req.SubscriptionCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)

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
	regex = "/api/subscription/search\\?([^/]+)$"
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

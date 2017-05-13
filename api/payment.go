package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/logger"
	"github.com/msproject/relive/util"
)

// PaymentAPI struct
type PaymentAPI struct {
	PaymentDBI        dbi.PaymentTblDBI
	PaymentHistoryDBI dbi.PaymentHistoryTblDBI
	LogObj            *logger.Logger
}

//	/api/payment/search
func handlePaymentSearch(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	// check for API Method
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/payment/search")
	}

	ID := args[1]
	idInt, errs := strconv.Atoi(ID)

	if errs != nil {
		http.Error(w, errs.Error(), http.StatusInternalServerError)
		return errs
	}

	var pays []util.PaymentDetails
	pays, err := api.PaymentDBI.SearchPayment(idInt)

	jsonStr, err := json.Marshal(pays)
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

// /api/payment/do
func handlePaymentDo(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {

	// check for API Method
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/payment/do")
	}

	// decode the JSON against the structure
	var req *dbmodel.PaymentEntry
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	err := api.PaymentDBI.AddPayment(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

// /api/payment/update -
func handlePaymentUpdate(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	// check for API Method
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/payment/update")
	}

	// decode the JSON against the structure
	var req *dbmodel.PaymentEntry
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	err := api.PaymentDBI.UpdatePayment(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// /api/payment/delete -
func handlePaymentDelete(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	// check for API Method
	if r.Method != "DELETE" {
		w.Header().Set("Allow", "DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/payment/delete")
	}

	// decode the JSON against the structure
	var req *dbmodel.PaymentEntry
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Error decoding the request: %s", err.Error())
	}

	if req.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required parameters NOT specified in delete request")
	}

	err := api.PaymentDBI.DeletePayment(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

type paymentT struct {
	regex string
	re    *regexp.Regexp
	f     func(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error
}

var payment []paymentT

func (api PaymentAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, d := range payment {
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
	regex = "/api/payment/search\\?([^/]+)$"
	//regex = "/api/payment/search$"
	payment = append(payment,
		paymentT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handlePaymentSearch,
		},
	)
	regex = "/api/payment/do$"
	payment = append(payment,
		paymentT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handlePaymentDo,
		},
	)
	regex = "/api/payment/update$"
	payment = append(payment,
		paymentT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handlePaymentUpdate,
		},
	)
	regex = "/api/payment/delete$"
	payment = append(payment,
		paymentT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handlePaymentDelete,
		},
	)
}

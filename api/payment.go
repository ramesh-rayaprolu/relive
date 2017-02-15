package api

import (
	"fmt"
	"net/http"
	"regexp"

	"../dbi"
	"../logger"
)

// PaymentAPI struct
type PaymentAPI struct {
	PaymentDBI        dbi.PaymentTblDBI
	PaymentHistoryDBI dbi.PaymentHistoryTblDBI
	LogObj            *logger.Logger
}

//	/api/payment/search
func handlePaymentSearch(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/payment/do
func handlePaymentDo(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// /api/payment/update -
func handlePaymentUpdate(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/payment/delete -
func handlePaymentDelete(api PaymentAPI, args []string, w http.ResponseWriter, r *http.Request) error {
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
	regex = "/api/payment/search$"
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

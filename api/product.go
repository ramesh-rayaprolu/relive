package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/logger"
	"github.com/msproject/relive/util"
)

// ProductsAPI struct
type ProductsAPI struct {
	ProductDBI dbi.ProductTblDBI
	LogObj     *logger.Logger
}

//InitProductsDB - create product if it doesnt exist
func InitProductsDB(sqlDBI dbi.DBI) (err error) {
	var exists bool

	exists, err = sqlDBI.ProductDBI.CheckProductTableExists()

	if err != nil {
		fmt.Printf("Error checking for Product Table < %s >\n", err.Error())
		return err
	}

	if !exists {
		fmt.Println("Product Table Does not exist")
		return nil
	}

	req := []util.CreateProductReq{{
		ProductID:   1001,
		ProductType: "bronze",
		StoreSize:   100,
		Duration:    30,
		Amount:      100,
	}, {
		ProductID:   1002,
		ProductType: "silver",
		StoreSize:   200,
		Duration:    30,
		Amount:      200,
	}, {
		ProductID:   1003,
		ProductType: "gold",
		StoreSize:   300,
		Duration:    30,
		Amount:      300,
	}}
	err = sqlDBI.ProductDBI.CreateProduct(req)
	if err != nil {
		fmt.Printf("Error creating product %s\n", err.Error())
		return err
	}
	return nil
}

func handleListProducts(api ProductsAPI, args []string, w http.ResponseWriter, r *http.Request) (err error) {

	var resp []dbmodel.ProductEntry

	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("Incorrect Method used for API /api/accounts/Search")
	}

	resp, err = api.ProductDBI.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err = writeResponse(resp, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}

type productT struct {
	regex string
	re    *regexp.Regexp
	f     func(api ProductsAPI, args []string, w http.ResponseWriter, r *http.Request) error
}

var product []productT

func (api ProductsAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, d := range product {
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
	regex = "/api/products/list$"
	product = append(product,
		productT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleListProducts,
		},
	)
}

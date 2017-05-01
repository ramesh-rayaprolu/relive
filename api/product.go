package api

import (
	//"encoding/json"
	"fmt"
	//"net/http"
	//"net/url"
	//"regexp"

	"github.com/msproject/relive/dbi"
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

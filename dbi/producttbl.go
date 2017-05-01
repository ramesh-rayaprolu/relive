package dbi

import (
	//"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/util"
)

// ProductTblDBI - testing
type ProductTblDBI interface {

	// CheckAccountTableExists - test
	CheckProductTableExists() (bool, error)

	// CreateProduct - create product
	CreateProduct(req []util.CreateProductReq) error
}

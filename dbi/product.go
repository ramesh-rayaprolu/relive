package dbi

import (
	"github.com/msproject/relive/dbmodel"
)

// ProductTblDBI - testing
type ProductTblDBI interface {
	// AddProduct - testing
	AddProduct(prDetails *dbmodel.ProductEntry) error
}

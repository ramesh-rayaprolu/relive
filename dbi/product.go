package dbi

import (
	"../dbmodel"
)

// ProductTblDBI - testing
type ProductTblDBI interface {
	// AddProduct - testing
	AddProduct(prDetails *dbmodel.ProductEntry) error
}

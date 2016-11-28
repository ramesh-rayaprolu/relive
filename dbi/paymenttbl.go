package dbi

import (
	"../dbmodel"
)

// PaymentTblDBI - testing
type PaymentTblDBI interface {
	// AddPayment - testing
	AddPayment(pyDetails *dbmodel.PaymentEntry) error
}

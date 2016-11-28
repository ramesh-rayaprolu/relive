package dbi

import (
	"github.com/msproject/relive/dbmodel"
)

// PaymentTblDBI - testing
type PaymentTblDBI interface {
	// AddPayment - testing
	AddPayment(pyDetails *dbmodel.PaymentEntry) error
}

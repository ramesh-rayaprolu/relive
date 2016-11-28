package dbi

import (
	"../dbmodel"
)

// PaymentHistoryTblDBI - testing
type PaymentHistoryTblDBI interface {
	// AddPaymentHistory - testing
	AddPaymentHistory(pyhDetails *dbmodel.PaymentHistoryEntry) error
}

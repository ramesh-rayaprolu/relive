package dbi

import (
	"github.com/msproject/relive/dbmodel"
)

// PaymentHistoryTblDBI - testing
type PaymentHistoryTblDBI interface {
	// AddPaymentHistory - testing
	AddPaymentHistory(pyhDetails *dbmodel.PaymentHistoryEntry) error
}

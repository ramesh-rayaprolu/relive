package dbi

import (
	"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/util"
)

// PaymentTblDBI - testing
type PaymentTblDBI interface {
	AddPayment(pyDetails *dbmodel.PaymentEntry) error
	SearchPayment(ID int) ([]util.PaymentDetails, error)
	UpdatePayment(pyDetails *dbmodel.PaymentEntry) error
	DeletePayment(paymentID int) error
}

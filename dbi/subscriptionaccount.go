package dbi

import (
	"github.com/msproject/relive/dbmodel"
)

// SubscriptionAccountTblDBI - testing
type SubscriptionAccountTblDBI interface {
	// AddSubscriptionAccount - testing
	AddSubscriptionAccount(subacDetails *dbmodel.SubscriptionAccountEntry) error
}

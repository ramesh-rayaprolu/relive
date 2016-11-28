package dbi

import (
	"../dbmodel"
)

// SubscriptionAccountTblDBI - testing
type SubscriptionAccountTblDBI interface {
	// AddSubscriptionAccount - testing
	AddSubscriptionAccount(subacDetails *dbmodel.SubscriptionAccountEntry) error
}

package dbi

import (
	"../dbmodel"
)

// SubscriptionTblDBI - testing
type SubscriptionTblDBI interface {
	// Addsubscription - testing
	AddSubscription(subDetails *dbmodel.SubscriptionEntry) error
}

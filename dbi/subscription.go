package dbi

import (
	"github.com/msproject/relive/dbmodel"
)

// SubscriptionTblDBI - testing
type SubscriptionTblDBI interface {
	// Addsubscription - testing
	AddSubscription(subDetails *dbmodel.SubscriptionEntry) error
}

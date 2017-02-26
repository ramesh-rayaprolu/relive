package dbi

import (
	//"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/util"
)

// SubscriptionTblDBI - testing
type SubscriptionTblDBI interface {
	CreateSubscription(req util.CreateSubscriptionReq) error
	UpdateSubscription(req util.CreateSubscriptionReq) error
	DeleteSubscription(subscriptionCode uint32) error
}

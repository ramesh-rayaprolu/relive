package dbi

import (
	"fmt"
	"sync"
	"time"
)

//DBI - testing
type DBI struct {
	AccountDBI             AccountTblDBI
	PaymentDBI             PaymentTblDBI
	PaymentHistoryDBI      PaymentHistoryTblDBI
	SubscriptionDBI        SubscriptionTblDBI
	SubscriptionAccountDBI SubscriptionAccountTblDBI
	ProductDBI             ProductTblDBI
	MediaTypeDBI           MediaTypeTblDBI
}

var dbi *DBI
var once sync.Once

// InitializeDBI - init
func InitializeDBI(svcAddr string, dbTimeout time.Duration) (DBI, error) {
	once.Do(func() {
		sqlDBI, sqlErr := NewMemSQLDBI(svcAddr, dbTimeout)
		if sqlErr != nil {
			return
		}
		dbi = &DBI{
			AccountDBI:             sqlDBI,
			PaymentDBI:             sqlDBI,
			PaymentHistoryDBI:      sqlDBI,
			SubscriptionDBI:        sqlDBI,
			SubscriptionAccountDBI: sqlDBI,
			ProductDBI:             sqlDBI,
			MediaTypeDBI:           sqlDBI,
		}
	})
	if dbi != nil {
		return *dbi, nil
	}
	return DBI{}, fmt.Errorf("DBI is not initialized")
}

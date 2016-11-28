package dbi

import (
	"../dbmodel"
)

// AccountTblDBI - testing
type AccountTblDBI interface {
	// AddAccounts - testing
	AddAccounts(acDetails *dbmodel.AccountEntry) error
}

package dbi

import (
	"../dbmodel"
	"../util"
)

// AccountTblDBI - testing
type AccountTblDBI interface {
	// CheckAccountTableExists - test
	CheckAccountTableExists() (bool, error)

	//CheckAccountExists - test
	CheckAccountExists(userName string) (bool, error)

	// Login - test
	Login(userName, PWD string) (*dbmodel.AccountEntry, error)

	// CreateAccount - test
	CreateAccount(req util.CreateAccountReq) error

	// AddAccounts - testing
	AddAccounts(acDetails *dbmodel.AccountEntry) error
}

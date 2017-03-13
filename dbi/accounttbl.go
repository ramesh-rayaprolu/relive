package dbi

import (
	"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/util"
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

	// SearchAccount - test
	SearchAccount(UserName string) (util.SearchAccountReq, error)

	// AddAccounts - testing
	AddAccounts(acDetails *dbmodel.AccountEntry) error
}

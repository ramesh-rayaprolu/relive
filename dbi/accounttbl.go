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

	//CheckAccountExists - test
	CheckAccountExistsByID(id uint64) error

	// Login - test
	Login(userName, PWD string) (*dbmodel.AccountEntry, error)

	// CreateAccount - test
	CreateAccount(req util.CreateAccountReq) error

	// SearchAccount - test
	SearchAccount(UserName string) (util.SearchAccountReq, error)

	UpdateAccount(upDetails *dbmodel.AccountEntry) error

	UpdateMyAccount(upDetails *dbmodel.AccountEntry) error

	// SearchAndGetAccountIDs - test
	SearchAndGetAccountIDs(adminID int) ([]util.UserDetails, error)

	// AddAccounts - testing
	AddAccounts(acDetails *dbmodel.AccountEntry) error

	DeleteAccount(userName string) error
}

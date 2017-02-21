package dbi

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/logger"
	"github.com/msproject/relive/util"
	"io"
	"math/rand"
	"time"
)

// SQLIF defines SQL database access functions
type SQLIF interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
}

// SQLDBI - testing
type SQLDBI struct {
	accessStr string
	db        SQLIF
	timeout   time.Duration
	logObj    *logger.Logger
}

// NewSQLDBI - testing
func NewSQLDBI(dsn string, timeout time.Duration, logObj *logger.Logger) (sqlDBI *SQLDBI, err error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	sqlDBI = &SQLDBI{
		accessStr: dsn,
		timeout:   timeout,
		db:        db,
		logObj:    logObj,
	}
	return //
}

//CheckAccountExists - check if given account exists
func (sqlDbi *SQLDBI) CheckAccountExists(userName string) (bool, error) {
	const IsAccountExistQuery = "Select COUNT(*) as count from Account where UserName = ?"
	var (
		rows *sql.Rows
		err  error
	)

	rows, err = sqlDbi.db.Query(IsAccountExistQuery, userName)
	if err != nil {
		sqlDbi.logObj.PrintError("Failed querying accounts %v", err)
		return false, fmt.Errorf("Failed querying accounts %v", err)
	}
	defer rows.Close()

	if checkCount(rows) > 0 {
		return true, nil
	}

	return false, nil
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0
		}
	}
	return count
}

//CheckAccountTableExists - check if accounts table exists
func (sqlDbi *SQLDBI) CheckAccountTableExists() (bool, error) {
	const checkTblQuery = `SHOW TABLES LIKE 'Account'`

	query := checkTblQuery
	args := []interface{}{}

	rows, err := sqlDbi.db.Query(query, args...)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	var tName []uint8
	if rows.Next() {
		err = rows.Scan(&tName)
		if err != nil {
			sqlDbi.logObj.PrintError("CheckAccountTableExists returns DB scan error: %s", err.Error())
			fmt.Printf("CheckAccountTableExists returns DB scan error: %s\n", err.Error())
			return false, err
		}
	}

	/* if there exists a table, tName will have it, else Account table does
	 *  not exist */
	if tName != nil {
		return true, nil
	}

	return false, nil
}

// CreateAccount - function to create an account row.
//                 Duplicate rows are not allowed and will throw error
func (sqlDbi *SQLDBI) CreateAccount(req util.CreateAccountReq) error {
	const createAccountQuery = `INSERT INTO Account (PID, UserName, FirstName, LastName, EmailID, PasswdDigest, Salt, Role) VALUES `
	var err error

	// Get passwordDigest and salt here
	passwordDigest, salt := saltedHash(req.PWD)

	query := createAccountQuery
	args := []interface{}{}

	query += "(?, ?, ?, ?, ?, ?, ?, ?)"

	args = append(args, req.CompanyID, req.UserName, req.FirstName, req.LastName, req.Email, passwordDigest, salt, req.Role)

	_, err = sqlDbi.db.Exec(query, args...)
	if err != nil {
		sqlDbi.logObj.PrintError("Failed to create account: %s", err.Error())
		return fmt.Errorf("Failed to create the account %v", err)
	}

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz@#$ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randStringBytes, saltedHash, hashPWDAndSalt functions:
//      - used to generate a random `salt`value and encrypt user password
func randStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func saltedHash(PWD string) (string, string) {
	h := md5.New()
	io.WriteString(h, PWD)
	pwmd5 := fmt.Sprintf("%x", h.Sum(nil))

	// Generate random salt
	salt := randStringBytes(10)

	// salt + MD5 splicing
	io.WriteString(h, salt)
	io.WriteString(h, pwmd5)
	passwordDigest := fmt.Sprintf("%x", h.Sum(nil))
	return passwordDigest, salt
}

func hashPWDAndSalt(PWD, salt string) string {
	h := md5.New()
	io.WriteString(h, PWD)
	pwmd5 := fmt.Sprintf("%x", h.Sum(nil))

	// salt + MD5 splicing
	io.WriteString(h, salt)
	io.WriteString(h, pwmd5)
	passwordDigest := fmt.Sprintf("%x", h.Sum(nil))
	return passwordDigest
}

//Login - verify user/password from DB
func (sqlDbi *SQLDBI) Login(userName, PWD string) (*dbmodel.AccountEntry, error) {
	/*    const LoginQuery = `SELECT UserName, Email, FullName, Role, Status, PasswordDigest, Salt
	        FROM Accounts WHERE UserName = ?`

	      var (
	          rows *sql.Rows
	          err  error
	      )

	      vals := []interface{}{}
	      vals = append(vals, userName)

	      rows, err = s.db.Query(LoginQuery, vals...)
	      if err != nil {
	          return nil, fmt.Errorf("Failed querying accounts %v", err)
	      }
	      defer rows.Close()
	      account := &model.Account{}
	      var passwordDigest, salt string
	      if rows.Next() {
	          err := rows.Scan(
	              &account.UserName,
	              &account.Email,
	              &account.FullName,
	              &account.Role,
	              &account.Status,
	              &passwordDigest,
	              &salt)
	          if err != nil {
	              return nil, fmt.Errorf("Failed scanning accounts %v", err)
	          }
	      }


	      inPasswordDigest := hashPWDAndSalt(PWD, salt)

	      if inPasswordDigest != passwordDigest {
	          return nil, fmt.Errorf("The username and password don't match.")
	      }

	      return account, nil*/
	return nil, nil
}

// AddAccounts - testing
func (sqlDbi *SQLDBI) AddAccounts(acDetails *dbmodel.AccountEntry) (err error) {

	const sqlInsertAccountQry = `INSERT INTO Account (ID, PID, FirstName, LastName, EmailID, PasswdDigest, Role) VALUES `
	//const sqlUpdateAccountQry = `ON DUPLICATE KEY UPDATE Name = VALUES(Name), Number = VALUES(Number) `

	var query = sqlInsertAccountQry
	args := []interface{}{}

	query += "(?, ?, ?, ?, ?, ?, ?)"
	args = append(args, acDetails.ID, acDetails.PID, acDetails.FirstName, acDetails.LastName, acDetails.EmailID, acDetails.PasswdDigest, acDetails.Role)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

// AddPayment - testing
func (sqlDbi *SQLDBI) AddPayment(pyDetails *dbmodel.PaymentEntry) (err error) {

	const sqlInsertPaymentQry = `INSERT INTO payment (ID, CCNumber, BillingAddress, CCExpiry, CVVCode) VALUES `

	var query = sqlInsertPaymentQry
	args := []interface{}{}

	query += "(?, ?, ?, ?, ?)"
	args = append(args, pyDetails.ID, pyDetails.CCNumber, pyDetails.BillingAddress, pyDetails.CCExpiry, pyDetails.CVVCode)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

// AddPaymentHistory - testing
func (sqlDbi *SQLDBI) AddPaymentHistory(pyhDetails *dbmodel.PaymentHistoryEntry) (err error) {

	const sqlInsertPaymenthistoryQry = `INSERT INTO Paymenthistory (ID, LastPaidState, LastType) VALUES `

	var query = sqlInsertPaymenthistoryQry
	args := []interface{}{}

	query += "(?, ?, ?)"
	args = append(args, pyhDetails.ID, pyhDetails.LastPaidState, pyhDetails.LastType)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

// CreateSubscription - testing
func (sqlDbi *SQLDBI) CreateSubscription(req util.CreateSubscriptionReq) (err error) {

	const sqlInsertSubscriptionQry = `INSERT INTO Subscription (ID, ProductID, ProductType, StoreLocation, StartDate, EndDate, NumberOfAdmins) VALUES `

	var query = sqlInsertSubscriptionQry
	args := []interface{}{}

	query += "(?, ?, ?, ?, ?, ?, ?)"
	args = append(args, req.ID, req.ProductID, req.ProductType, req.StoreLocation, req.StartDate, req.EndDate, req.NumberOfAdmins)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

// AddSubscriptionAccount - testing
func (sqlDbi *SQLDBI) AddSubscriptionAccount(subacDetails *dbmodel.SubscriptionAccountEntry) (err error) {

	const sqlInsertSubscriptionaccountQry = `INSERT INTO Subscriptionaccount (ID, PID) VALUES `

	var query = sqlInsertSubscriptionaccountQry
	args := []interface{}{}

	query += "(?, ?)"
	args = append(args, subacDetails.ID, subacDetails.PID)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

// AddProduct - testing
func (sqlDbi *SQLDBI) AddProduct(prDetails *dbmodel.ProductEntry) (err error) {

	const sqlInsertProductsQry = `INSERT INTO Product (ProductID, ProductType, StoreSize, Duration, Amount) VALUES `

	var query = sqlInsertProductsQry
	args := []interface{}{}

	query += "(?, ?, ?, ?, ?)"
	args = append(args, prDetails.ProductID, prDetails.ProductType, prDetails.StoreSize, prDetails.Duration, prDetails.Amount)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

// AddMediaType - testing
func (sqlDbi *SQLDBI) AddMediaType(mtDetails *dbmodel.MediaTypeEntry) (err error) {

	const sqlInsertMediatypeQry = `INSERT INTO Mediatype (ID, Catalog, FileName, Title, Description, URL, Poster) VALUES `

	var query = sqlInsertMediatypeQry
	args := []interface{}{}

	query += "(?, ?, ?, ?, ?, ?, ?)"
	args = append(args, mtDetails.ID, mtDetails.Catalog, mtDetails.FileName, mtDetails.Title, mtDetails.Description, mtDetails.URL, mtDetails.Poster)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

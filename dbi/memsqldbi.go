package dbi

import (
	"../dbmodel"
	"database/sql"
	"time"
)

// SQLIF defines SQL database access functions
type SQLIF interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
}

// MemSQLDBI - testing
type MemSQLDBI struct {
	accessStr string
	db        SQLIF
	timeout   time.Duration
}

// NewMemSQLDBI - testing
func NewMemSQLDBI(dsn string, timeout time.Duration) (sqlDBI *MemSQLDBI, err error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	sqlDBI = &MemSQLDBI{
		accessStr: dsn,
		timeout:   timeout,
		db:        db,
	}
	return //
}

// AddAccounts - testing
func (sqlDbi *MemSQLDBI) AddAccounts(acDetails *dbmodel.AccountEntry) (err error) {

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
func (sqlDbi *MemSQLDBI) AddPayment(pyDetails *dbmodel.PaymentEntry) (err error) {

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
func (sqlDbi *MemSQLDBI) AddPaymentHistory(pyhDetails *dbmodel.PaymentHistoryEntry) (err error) {

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

// AddSubscription - testing
func (sqlDbi *MemSQLDBI) AddSubscription(subDetails *dbmodel.SubscriptionEntry) (err error) {

	const sqlInsertSubscriptionQry = `INSERT INTO Subscription (ID, ProductID, ProductType, StoreLocation, StartDate, EndDate, NumberOfAdmins) VALUES `

	var query = sqlInsertSubscriptionQry
	args := []interface{}{}

	query += "(?, ?, ?, ?, ?, ?, ?)"
	args = append(args, subDetails.ID, subDetails.ProductID, subDetails.ProductType, subDetails.StoreLocation, subDetails.StartDate, subDetails.EndDate, subDetails.NumberOfAdmins)
	//query += sqlUpdateAccountQry

	_, err = sqlDbi.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

// AddSubscriptionAccount - testing
func (sqlDbi *MemSQLDBI) AddSubscriptionAccount(subacDetails *dbmodel.SubscriptionAccountEntry) (err error) {

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
func (sqlDbi *MemSQLDBI) AddProduct(prDetails *dbmodel.ProductEntry) (err error) {

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
func (sqlDbi *MemSQLDBI) AddMediaType(mtDetails *dbmodel.MediaTypeEntry) (err error) {

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

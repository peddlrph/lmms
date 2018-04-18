// Package smartmoney
package smartmoney

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table           = "smartmoney"
	smartmoneytable = "smartmoney"
	transtable      = "transaction"
	cashtable       = "cash"
)

// Item defines the model.
type Item struct {
	ID                       uint32         `db:"id"`
	Trans_Datetime           mysql.NullTime `db:"trans_datetime"`
	Trans_Datetime_Formatted string
	Amount                   sql.NullFloat64 `db:"amount"`
	Amount_String            string
	Details                  string
	Details_Split            []string
	CreatedAt                mysql.NullTime `db:"created_at"`
	UpdatedAt                mysql.NullTime `db:"updated_at"`
	DeletedAt                mysql.NullTime `db:"deleted_at"`
}

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// ByID gets an item by ID.
func ByID(db Connection, ID string) (Item, bool, error) {
	result := Item{}
	err := db.Get(&result, fmt.Sprintf(`
		SELECT id, trans_datetime, amount, details, created_at, updated_at, deleted_at
		FROM %v
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID)
	return result, err == sql.ErrNoRows, err
}

// All gets all items.
func All(db Connection) ([]Item, float32, bool, error) {
	var result []Item
	var sum float32
	err := db.Select(&result, fmt.Sprintf(`
		SELECT id, trans_datetime, amount, details, created_at, updated_at, deleted_at
		FROM %v
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		`, table))
	_ = db.Get(&sum, fmt.Sprintf(`
		SELECT sum(amount)
		FROM %v
		WHERE deleted_at IS NULL
		LIMIT 1
		`, table))
	return result, sum, err == sql.ErrNoRows, err
}

// Create adds an item.
func Create(db Connection, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(now(),?,?)
		`, table),
		amt, details)
	return result, err
}

func SmartMoneyWithCash(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	trans_code := "SmartMoneyWithCash"

	trans_details := "SmartMoneyWithCash: |"
	trans_details = trans_details + "-  Subtract " + amount + " from cash.|"
	trans_details = trans_details + "-  Add " + amount + " to smartmoney.|"
	trans_details = trans_details + "Details: " + details
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, transtable), trans_date,
		amt, trans_details)
	if err != nil {
		return result, err
	}
	trans_id, _ := result.LastInsertId()
	transactiontag := " Trans#: " + strconv.FormatInt(trans_id, 10)

	cash_details := "SmartMoneyWithCash: |"
	cash_details = cash_details + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt*(-1), cash_details)
	if err != nil {
		return result, err
	}

	sm_details := "SmartMoneyWithCash: |"
	sm_details = sm_details + "Details: " + details + "|"
	sm_details = sm_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, smartmoneytable), trans_date, trans_code,
		amt, sm_details)
	return result, err
}

func EncashSmartMoney(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := fmt.Sprintf("%.2f", amt*(0.005))

	trans_code := "EncashSmartMoney"

	trans_details := "EncashSmartMoney: |"
	trans_details = trans_details + "-  Add " + amount + " to smartmoney.|"
	trans_details = trans_details + "-  Subtract " + amount + " from cash.|"
	trans_details = trans_details + "-  Add " + fee + " to cash as fee.|"
	trans_details = trans_details + "Details: " + details
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, transtable), trans_date,
		amt, trans_details)
	if err != nil {
		return result, err
	}
	trans_id, _ := result.LastInsertId()
	transactiontag := " Trans#: " + strconv.FormatInt(trans_id, 10)

	cash_details := "EncashSmartMoney: |"
	cash_details = cash_details + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt*(-1), cash_details)
	if err != nil {
		return result, err
	}

	//cash_details = "EncashSmartMoney: |"
	//cash_details = cash_details + details + "|"
	//cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,fee,details)
		VALUES
		(?,?,?,?,?)
		`, cashtable), trans_date, trans_code,fee,1,cash_details)
	if err != nil {
		return result, err
	}

	sm_details := "EncashSmartMoney: |"
	sm_details = sm_details + "Details: " + details + "|"
	sm_details = sm_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, smartmoneytable), trans_date, trans_code,
		amt, sm_details)
	return result, err
}

//func Create(db Connection, user_id string, customer string, order_detail string, php_price string, bch_price float64, pricephp float64, order_validity int) (sql.Result, error) {
//	result, err := db.Exec(fmt.Sprintf(`
//		INSERT INTO %v
//		(user_id,customer,order_detail,php_price,bch_price, exchange_rate,created_at,expired_at)
//		VALUES
//		(?,?,?,?,?,?,now(),date_add(now(), interval ? minute))
//		`, table),
//		user_id, customer, order_detail, php_price, bch_price, pricephp, order_validity)
//	return result, err
//}

// Update makes changes to an existing item.
func Update(db Connection, trans_datetime string, amount string, details string, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET trans_datetime=?, amount = ?,details=?
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		trans_datetime, amount, details, ID)
	return result, err
}

// DeleteHard removes an item.
func DeleteHard(db Connection, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		DELETE FROM %v
		WHERE id = ?
			AND deleted_at IS NULL
		`, table),
		ID)
	return result, err
}

// DeleteSoft marks an item as removed.
func DeleteSoft(db Connection, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET deleted_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID)
	return result, err
}

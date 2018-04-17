// Package smartpadala
package smartpadala

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table            = "smartpadala"
	cashtable        = "cash"
	loadstable       = "loads"
	smartmoneytable  = "smartmoney"
	transtable       = "transaction"
	smartpadalatable = "smartpadala"
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
		SELECT id, name, created_at, updated_at, deleted_at
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
	if err != nil {
		fmt.Println(err)
	}
	return result, sum, err == sql.ErrNoRows, err
}

// Create adds an item.
func Create(db Connection, name string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(name)
		VALUES
		(?)
		`, table),
		name)
	return result, err
}

func ReceiveSP(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := get_fee(amt)
	fee_string := fmt.Sprintf("%.2f", fee)
	trans_details := "ReceiveSP: |"
	trans_details = trans_details + "-  Subtract " + amount + " from cash.|"
	trans_details = trans_details + "-  Add " + amount + " to smartmoney.|"
	trans_details = trans_details + "-  Add " + fee_string + " to smartmoney as fee.|"
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

	cash_details := "ReceiveSP: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, cashtable), trans_date,
		amt*(-1), cash_details)

	sm_details := "ReceiveSP: |"
	sm_details = sm_details + "Details: " + details + "|"
	sm_details = sm_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartmoneytable), trans_date,
		amt, sm_details)
	if err != nil {
		return result, err
	}

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartmoneytable), trans_date,
		fee, sm_details)

	sp_details := "ReceiveSP: |"
	sp_details = sp_details + "Details: " + details + "|"
	sp_details = sp_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartpadalatable), trans_date,
		amt, sp_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartpadalatable), trans_date,
		fee, sp_details)

	if err != nil {
		return result, err
	}
	return result, err
}

func ReceiveSPX(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)
	trans_details := details + " ReceiveSP: "
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, cashtable), trans_date,
		amt*(-1), details)
	trans_details = trans_details + "Subtract " + amount + " from cash."
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartmoneytable), trans_date,
		amt, details)
	trans_details = trans_details + " Add " + amount + " to smartmoney."
	//fee_rate := get_fee_rate("receivesp",amt)
	fee := fmt.Sprintf("%.2f", amt*(0.03))
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartmoneytable), trans_date,
		fee, details)
	trans_details = trans_details + " Add " + fee + " to smartmoney as fee."
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, transtable), trans_date,
		amt, trans_details)
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, table), trans_date,
		amt, trans_details)
	return result, err
}

func SendSP(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := get_fee(amt)
	fee_string := fmt.Sprintf("%.2f", fee)
	trans_details := "SendSP: |"
	trans_details = trans_details + "-  Subtract " + amount + " from smartmoney.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee_string + " to smartmoney as fee.|"
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

	cash_details := "SendSP: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, cashtable), trans_date,
		amt, cash_details)

	sm_details := "SendSP: |"
	sm_details = sm_details + "Details: " + details + "|"
	sm_details = sm_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartmoneytable), trans_date,
		amt*(-1), sm_details)
	if err != nil {
		return result, err
	}

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartmoneytable), trans_date,
		fee, sm_details)

	sp_details := "SendSP: |"
	sp_details = sp_details + "Details: " + details + "|"
	sp_details = sp_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartpadalatable), trans_date,
		amt, sp_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartpadalatable), trans_date,
		fee, sp_details)

	if err != nil {
		return result, err
	}
	return result, err
}

func SendSPX(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)
	trans_details := details + " SendSP: "
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, cashtable), trans_date,
		amt, details)
	trans_details = trans_details + "Add " + amount + " to cash."
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, smartmoneytable), trans_date,
		amt*(-1), details)
	trans_details = trans_details + " Subtract " + amount + " from smartmoney."
	fee_rate := get_fee_rate("ReceiveSP")
	fee := fmt.Sprintf("%.2f", amt*fee_rate)
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, cashtable), trans_date,
		fee, details)
	trans_details = trans_details + " Add " + fee + " to cash as fee."
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, transtable), trans_date,
		amt, trans_details)
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, table), trans_date,
		amt, trans_details)
	return result, err
}

// Update makes changes to an existing item.
func Update(db Connection, name string, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET name = ?
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		name, ID)
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

func get_fee_rate(transcode string) float64 {
	switch transcode {
	case "SendSP":
		return 0.01
	case "ReceiveSP":
		return 0.02
	default:
		return 0.01
	}
}

func get_fee(amt float64) float64 {
	switch true {
	case amt <= 1000:
		return 11.00
	case amt > 1000:
		return 20.00
	default:
		return 1000
	}
}

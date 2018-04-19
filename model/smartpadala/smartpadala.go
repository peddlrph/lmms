// Package smartpadala
package smartpadala

import (
	"database/sql"
	"fmt"
	"math"
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

	fee := get_receivefee(amt, 1000, 11.5)
	fee_string := fmt.Sprintf("%.2f", fee)

	trans_code := "ReceiveSP"

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

	sm_details := "ReceiveSP: |"
	sm_details = sm_details + "Details: " + details + "|"
	sm_details = sm_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, smartmoneytable), trans_date, trans_code,
		amt, sm_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,fee,details)
		VALUES
		(?,?,?,?,?)
			`, smartmoneytable), trans_date, trans_code,
		fee, 1, sm_details)

	cash_details := "ReceiveSP: |"
	cash_details = cash_details + "Details: " + details + "|"
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

	/*
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
	*/
	return result, err
}

func SendSP(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	_, fee2sender, fee2receiver := get_sendfees(amt)

	//fee_string := fmt.Sprintf("%.2f", fee)
	fee2sender_string := fmt.Sprintf("%.2f", fee2sender)
	fee2receiver_string := fmt.Sprintf("%.2f", fee2receiver)

	trans_code := "SendSP"

	trans_details := "SendSP: |"
	trans_details = trans_details + "-  Subtract " + amount + " from smartmoney.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee2sender_string + " to cash as fee.|"
	trans_details = trans_details + "-  Subtract " + fee2receiver_string + " from smartmoney. Fee to Smart and Receiver.|"
	trans_details = trans_details + "-  Add " + fee2receiver_string + " to cash. Fee to Smart and Receiver.|"
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
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt, cash_details)

	sm_details := "SendSP: |"
	sm_details = sm_details + "Details: " + details + "|"
	sm_details = sm_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, smartmoneytable), trans_date, trans_code,
		amt*(-1), sm_details)
	if err != nil {
		return result, err
	}

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,fee,details)
		VALUES
		(?,?,?,?,?)
		`, cashtable), trans_date, trans_code,
		fee2sender, 1, cash_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		fee2receiver, cash_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, smartmoneytable), trans_date, trans_code,
		fee2receiver*(-1), sm_details)

	/*
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

	*/

	if err != nil {
		return result, err
	}
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

func get_receivefee(amt float64, unit float64, per_unit_rate float64) float64 {

	fee_cnt := math.Ceil(amt / unit)

	//fmt.Println(math.Ceil(fee_cnt))
	//fmt.Println(math.Ceil(fee_cnt))
	//fmt.Println(math.Mod(amt, unit))

	//if math.Mod(amt, unit) > 0 {
	return fee_cnt * per_unit_rate
	//} else {
	//	return fee_cnt * per_unit_rate
	//}

}

func get_sendfees(amt float64) (float64, float64, float64) {

	fee_cnt := math.Ceil(amt / 1000)
	fee := fee_cnt * 30
	fee2sender := fee_cnt * 9.5
	fee2receiver := fee_cnt * (30 - 9.5)

	//fmt.Println(math.Ceil(fee_cnt))
	//fmt.Println(math.Ceil(fee_cnt))
	//fmt.Println(math.Mod(amt, unit))

	//if math.Mod(amt, unit) > 0 {
	return fee, fee2sender, fee2receiver
	//} else {
	//	return fee_cnt * per_unit_rate
	//}

}

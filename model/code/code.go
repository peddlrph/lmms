// Package code
package code

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	//"github.com/peddlrph/lib/utilities"
)

var (
	// table is the table name.
	table      = "codes"
	codestable = "codes"
	cashtable  = "cash"
	transtable = "transaction"
)

// Item defines the model.
type Item struct {
	ID                       uint32         `db:"id"`
	Trans_Datetime           mysql.NullTime `db:"trans_datetime"`
	Trans_Datetime_Formatted string
	Code_Count               int32           `db:"code_count"`
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
		SELECT id, trans_datetime, code_count,amount, details, created_at, updated_at, deleted_at
		FROM %v
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID)
	return result, err == sql.ErrNoRows, err
}

// All gets all items.
func All(db Connection) ([]Item, float32, float32, bool, error) {
	var result []Item
	var sum, cnt float32
	err := db.Select(&result, fmt.Sprintf(`
		SELECT id, trans_datetime, code_count,amount, details, created_at, updated_at, deleted_at
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
	_ = db.Get(&cnt, fmt.Sprintf(`
		SELECT sum(code_count)
		FROM %v
		WHERE deleted_at IS NULL
		LIMIT 1
		`, table))
	return result, sum, cnt, err == sql.ErrNoRows, err
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

func BuyCode(db Connection, trans_date string, codecount string, details string) (sql.Result, error) {
	//cost_per_code := uint64(3698)
	//amt
	cnt, _ := strconv.ParseInt(codecount, 10, 64)
	amount := strconv.FormatInt(3698*cnt, 10)
	amt, _ := strconv.ParseFloat(amount, 64)

	trans_code := "BuyCode"

	trans_details := "BuyCode: |"
	trans_details = trans_details + "-  Subtract " + amount + " from cash.|"
	trans_details = trans_details + "-  Add " + amount + " to codes.|"
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

	cash_details := "BuyCode: |"
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

	code_details := "BuyCode: |"
	code_details = code_details + "Details: " + details + "|"
	code_details = code_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,code_count,amount,details)
		VALUES
		(?,?,?,?,?)
		`, codestable), trans_date, trans_code, codecount,
		amt, code_details)
	return result, err
}

func Code2Muni(db Connection, trans_date string, codecount string, details string) (sql.Result, error) {
	//amt, _ := strconv.ParseFloat(amount, 64)

	cnt, _ := strconv.ParseInt(codecount, 10, 64)
	amount := strconv.FormatInt(3698*cnt, 10)
	amt, _ := strconv.ParseFloat(amount, 64)

	//fmt.Println(cnt)

	fee := 100.00 * cnt

	//fmt.Printf("%T : %v \n", fee, fee)
	trans_code := "Code2Muni"

	fee_string := fmt.Sprintf("%.2f", float64(fee))
	trans_details := "Code2Muni: |"
	trans_details = trans_details + "-  Subtract " + amount + " from code.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee_string + " to cash as fee.|"
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

	code_details := "Code2Muni: |"
	code_details = code_details + "Details: " + details + "|"
	code_details = code_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,code_count,amount,details)
		VALUES
		(?,?,?,?,?)
		`, codestable), trans_date,
		trans_code, cnt*(-1), amt*(-1), code_details)
	if err != nil {
		return result, err
	}

	cash_details := "Code2Muni: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt, cash_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		fee, cash_details)
	if err != nil {
		return result, err
	}
	return result, err
}

func Code2Brgy(db Connection, trans_date string, codecount string, details string) (sql.Result, error) {
	//amt, _ := strconv.ParseFloat(amount, 64)

	cnt, _ := strconv.ParseInt(codecount, 10, 64)
	amount := strconv.FormatInt(3698*cnt, 10)
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := 200.00 * cnt
	trans_code := "Code2Brgy"
	fee_string := fmt.Sprintf("%.2f", float64(fee))

	trans_details := "Code2Brgy: |"
	trans_details = trans_details + "-  Subtract " + amount + " from load.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee_string + " to cash as fee.|"
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

	code_details := "Code2Brgy: |"
	code_details = code_details + "Details: " + details + "|"
	code_details = code_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,code_count,amount,details)
		VALUES
		(?,?,?,?,?)
		`, codestable), trans_date, trans_code, cnt*(-1),
		amt*(-1), code_details)
	if err != nil {
		return result, err
	}

	cash_details := "Code2Brgy: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt, cash_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		fee, cash_details)
	if err != nil {
		return result, err
	}
	return result, err
}

func Code2Dealer(db Connection, trans_date string, codecount string, details string) (sql.Result, error) {
	//amt, _ := strconv.ParseFloat(amount, 64)

	cnt, _ := strconv.ParseInt(codecount, 10, 64)
	amount := strconv.FormatInt(3698*cnt, 10)
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := 300.00 * cnt
	trans_code := "Code2Dealer"
	fee_string := fmt.Sprintf("%.2f", float64(fee))

	trans_details := "Code2Dealer: |"
	trans_details = trans_details + "-  Subtract " + amount + " from load.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee_string + " to cash as fee.|"
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

	code_details := "Code2Dealer: |"
	code_details = code_details + "Details: " + details + "|"
	code_details = code_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,code_count,amount,details)
		VALUES
		(?,?,?,?,?)
		`, codestable), trans_date, trans_code, cnt*(-1),
		amt*(-1), code_details)
	if err != nil {
		return result, err
	}

	cash_details := "Code2Dealer: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt, cash_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		fee, cash_details)
	if err != nil {
		return result, err
	}
	return result, err
}

func Code2DealerY(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	cnt := 0

	fee := 300.00 * cnt
	fee_string := fmt.Sprintf("%.2f", float64(fee))
	trans_details := "Code2Dealer: |"
	trans_details = trans_details + "-  Subtract " + amount + " from load.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee_string + " to cash as fee.|"
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

	code_details := "Code2Dealer: |"
	code_details = code_details + "Details: " + details + "|"
	code_details = code_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, codestable), trans_date,
		amt*(-1), code_details)
	if err != nil {
		return result, err
	}

	cash_details := "Code2Dealer: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, cashtable), trans_date,
		amt, cash_details)

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, cashtable), trans_date,
		fee, cash_details)
	if err != nil {
		return result, err
	}
	return result, err
}

func Code2DealerX(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)
	trans_details := details + " Code2Dealer: "
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
		`, codestable), trans_date,
		amt*(-1), details)
	trans_details = trans_details + " Subtract " + amount + " from codes."
	//fee_rate := get_fee_rate("ReceiveSP")
	fee := fmt.Sprintf("%.2f", 300.00)
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
		`, table), trans_date,
		amt, trans_details)
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

// Package load
package load

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table           = "loads"
	loadstable      = "loads"
	cashtable       = "cash"
	smartmoneytable = "smartmoney"
	transtable      = "transaction"
)

// Item defines the model.
type Item struct {
	ID                       uint32         `db:"id"`
	Trans_Datetime           mysql.NullTime `db:"trans_datetime"`
	Trans_Datetime_Formatted string
	Trans_Code               sql.NullString  `db:"trans_code"`
	Mobile_Number            sql.NullString  `db:"mobile_number"`
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
		SELECT id, trans_datetime, trans_code, mobile_number, amount, details, created_at, updated_at, deleted_at
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
		SELECT id, trans_datetime, trans_code,mobile_number,amount, details, created_at, updated_at, deleted_at
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
		(trans_datetime,trans_code,mobile_number,amount,details)
		VALUES
		(now(),?,?)
		`, table),
		amt, details)
	return result, err
}

func ReplenishWithCash(db Connection, trans_date string, mobile_num string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)
	trans_code := "ReplenishWithCash"
	trans_details := "ReplenishWithCash: |"
	trans_details = trans_details + "-  Subtract " + amount + " from cash.|"
	trans_details = trans_details + "-  Add " + amount + " to load.|"
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

	cash_details := "ReplenishWithCash: |"
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

	loads_details := "ReplenishWithCash: |"
	loads_details = loads_details + "Details: " + details + "|"
	loads_details = loads_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,mobile_number,amount,details)
		VALUES
		(?,?,?,?,?)
		`, loadstable), trans_date, trans_code, mobile_num,
		amt, loads_details)
	return result, err
}

func ReplenishWithSM(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)
	trans_code := "ReplenishWithSM"
	trans_details := "ReplenishWithSM: |"
	trans_details = trans_details + "-  Subtract " + amount + " from smartmoney.|"
	trans_details = trans_details + "-  Add " + amount + " to load.|"
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

	sm_details := "ReplenishWithSM: |"
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

	loads_details := "ReplenishWithSM: |"
	loads_details = loads_details + "Details: " + details + "|"
	loads_details = loads_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, loadstable), trans_date, trans_code,
		amt, loads_details)
	return result, err
}

func Load2Muni(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := fmt.Sprintf("%.2f", amt*(0.01))

	trans_code := "Load2Muni"
	trans_details := "Load2Muni: |"
	trans_details = trans_details + "-  Subtract " + amount + " from load.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee + " to load as fee.|"
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

	loads_details := "Load2Muni: |"
	loads_details = loads_details + "Details: " + details + "|"
	loads_details = loads_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, loadstable), trans_date, trans_code,
		amt*(-1), loads_details)
	if err != nil {
		return result, err
	}

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,fee,details)
		VALUES
		(?,?,?,?,?)
		`, loadstable), trans_date, trans_code,
		fee, 1,loads_details)
	if err != nil {
		return result, err
	}

	cash_details := "Loads2Muni: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt, cash_details)
	return result, err
}

func Load2Brgy(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := fmt.Sprintf("%.2f", amt*(0.0125))

	trans_code := "Load2Brgy"
	trans_details := "Load2Brgy: |"
	trans_details = trans_details + "-  Subtract " + amount + " from load.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee + " to load as fee.|"
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

	loads_details := "Load2Brgy: |"
	loads_details = loads_details + "Details: " + details + "|"
	loads_details = loads_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, loadstable), trans_date, trans_code,
		amt*(-1), loads_details)
	if err != nil {
		return result, err
	}

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,fee,details)
		VALUES
		(?,?,?,?,?)
		`, loadstable), trans_date, trans_code,
		fee,1, loads_details)
	if err != nil {
		return result, err
	}

	cash_details := "Loads2Brgy: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt, cash_details)
	return result, err
}

func Load2Dealer(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)

	fee := fmt.Sprintf("%.2f", amt*(0.03))

	trans_code := "Load2Dealer"
	trans_details := "Load2Dealer: |"
	trans_details = trans_details + "-  Subtract " + amount + " from load.|"
	trans_details = trans_details + "-  Add " + amount + " to cash.|"
	trans_details = trans_details + "-  Add " + fee + " to load as fee.|"
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

	loads_details := "Load2Dealer: |"
	loads_details = loads_details + "Details: " + details + "|"
	loads_details = loads_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, loadstable), trans_date, trans_code,
		amt*(-1), loads_details)
	if err != nil {
		return result, err
	}

	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,fee,details)
		VALUES
		(?,?,?,?,?)
		`, loadstable), trans_date, trans_code,
		fee, 1,loads_details)
	if err != nil {
		return result, err
	}

	cash_details := "Loads2Dealer: |"
	cash_details = cash_details + "Details: " + details + "|"
	cash_details = cash_details + transactiontag
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,trans_code,amount,details)
		VALUES
		(?,?,?,?)
		`, cashtable), trans_date, trans_code,
		amt, cash_details)
	return result, err
}

func Load2DealerX(db Connection, trans_date string, amount string, details string) (sql.Result, error) {
	amt, _ := strconv.ParseFloat(amount, 64)
	trans_details := "Load2Dealer: "
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
		`, loadstable), trans_date,
		amt*(-1), details)
	trans_details = trans_details + " Subtract " + amount + " from load."
	fee := fmt.Sprintf("%.2f", amt*(0.03))
	result, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(trans_datetime,amount,details)
		VALUES
		(?,?,?)
		`, loadstable), trans_date,
		fee, details)
	trans_details = trans_details + " Add " + fee + " to load as fee."
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

// Package summary
package summary

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "summary"
)

// Item defines the model.
type Item struct {
	ID                     uint32         `db:"id"`
	SnapshotDate           mysql.NullTime `db:"snapshotdate"`
	SnapshotDate_Formatted string
	Cash                   sql.NullFloat64 `db:"cash"`
	Cash_String            string
	Loads                  sql.NullFloat64 `db:"loads"`
	Loads_String           string
	SmartMoney             sql.NullFloat64 `db:"smartmoney"`
	SmartMoney_String      string
	Codes                  sql.NullFloat64 `db:"codes"`
	Codes_String           string
	Total                  sql.NullFloat64 `db:"total"`
	Total_String           string
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
func All(db Connection) ([]Item, bool, error) {
	var result []Item

	_, err := db.Exec(fmt.Sprintf(`
		TRUNCATE TABLE %v
		`, table))

	//fmt.Println(err)
	//if err != nil {
	//	return _, err == sql.ErrNoRows, err
	//}

	_, err = db.Exec(`insert into summary(snapshotdate,cash,loads,smartmoney,codes,total) select curdate(), ifnull((select sum(amount) from cash),0),ifnull((select sum(amount) from loads),0),ifnull((select sum(amount) from smartmoney),0), ifnull((select sum(amount) from codes),0),ifnull((select sum(amount) from cash),0)+ifnull((select sum(amount) from loads),0)+ifnull((select sum(amount) from smartmoney),0)+ifnull((select sum(amount) from codes),0)  as total from dual`)
	//	if err != nil {
	//		return res, err == sql.ErrNoRows, err
	//	}
	//fmt.Println(err)

	err = db.Select(&result, fmt.Sprintf(`
		SELECT id, snapshotdate,cash,loads,smartmoney,codes,total
		FROM %v
		`, table))
	return result, err == sql.ErrNoRows, err
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

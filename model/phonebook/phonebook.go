// Package phonebook
package phonebook

import (
	"database/sql"
	"fmt"
	//	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "phonebook"
)

// Item defines the model.
type Item struct {
	MobileNumber string         `db:"mobile_number"`
	Name         sql.NullString `db:"name"`
	Category     sql.NullString `db:"category"`
	Location     sql.NullString `db:"location"`
}

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// ByID gets an item by ID.
func ByMobileNumber(db Connection, MobileNumber string) (Item, bool, error) {
	result := Item{}
	err := db.Get(&result, fmt.Sprintf(`
		SELECT mobile_number, name, category,location
		FROM %v
		WHERE mobile_number = ?
		LIMIT 1
		`, table),
		MobileNumber)
	return result, err == sql.ErrNoRows, err
}

// All gets all items.
func All(db Connection) ([]Item, bool, error) {
	var result []Item
	err := db.Select(&result, fmt.Sprintf(`
		SELECT mobile_number, name, category,location
		FROM %v
		`, table))
	return result, err == sql.ErrNoRows, err
}

func AddMobileNumberToPhonebook(db Connection, mobile_num string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT IGNORE INTO %v
		(mobile_number)
		VALUES
		(?)
		`, table),
		mobile_num)
	return result, err
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
func Update(db Connection, mobile_number string, name string, category string, location string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET name = ?, category = ?, location = ?
		WHERE mobile_number = ?
		LIMIT 1
		`, table),
		name, category, location, mobile_number)
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

// Package mobile_ip
package mobile_ip

import (
	"database/sql"
	"fmt"
	//"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "mobile_ip"
)

// Item defines the model.
type Item struct {
	ID         uint32 `db:"id"`
	IP_Address string `db:"ip_address"`
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
		SELECT id, ip_address
		FROM %v
		WHERE id = ?
		LIMIT 1
		`, table),
		ID)
	return result, err == sql.ErrNoRows, err
}

//func GetIP()(string,error){

//}

// All gets all items.
func All(db Connection) ([]Item, bool, error) {
	var result []Item
	err := db.Select(&result, fmt.Sprintf(`
		SELECT id, name, created_at, updated_at, deleted_at
		FROM %v
		WHERE deleted_at IS NULL
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
func Update(db Connection, ip_address string, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET ip_address = ?
		WHERE id = ?
		LIMIT 1
		`, table),
		ip_address, ID)
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

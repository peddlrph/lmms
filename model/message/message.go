// Package message
package message

import (
	"database/sql"
	"encoding/json"
	"fmt"
	//"net/http"
	//"github.com/go-sql-driver/mysql"
	"github.com/blue-jay/blueprint/model/mobile_ip"
	"github.com/peddlrph/lib/smsgateway"
	s "strings"
)

var (
	// table is the table name.
	table = "message"
)

// Item defines the model.
type Message struct {
	Address string `json:"address"`
	Body    string `json:"body"`
	Msg_Box string `json:"msg_box"`
	Id      string `json:"_id"`
}

type Item struct {
	Limit    string `json:"limit"`
	Offset   string `json:"offset"`
	Size     string `json:"size"`
	Messages []Message
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
func All(db Connection) ([]Message, Message, error) {
	// var result []Message
	//c := flight.Context(w, r)

	Balance_Message := Message{}

	Balance_Message.Body = "Balance not available"

	mobile_ip, _, _ := mobile_ip.ByID(db, "1")

	response, err := smsgateway.GetMessages("http://" + mobile_ip.IP_Address + ":8080/v1/sms/?limit=10000")
	if err == nil {
		//fmt.Println("yeah")
		smsgateway.WriteMessagesToFile(response)
		if err != nil {
			fmt.Println(err)
		}
	}

	mesgs := Item{}
	json.Unmarshal([]byte(response), &mesgs)

	//err := db.Select(&result, fmt.Sprintf(`
	//	SELECT id, name, created_at, updated_at, deleted_at
	//	FROM %v
	//	WHERE deleted_at IS NULL
	//	`, table))

	messages := mesgs.Messages

	strList := []string{"Load Wallet", "Commission Fund", "Smart Money"}

	for i := 0; i < len(messages); i++ {
		if s.Contains(messages[i].Body, strList[0]) && s.Contains(messages[i].Body, strList[1]) && s.Contains(messages[i].Body, strList[2]) {
			//fmt.Println(messages[i].Body)
			Balance_Message = messages[i]
			break
		}
	}

	return messages, Balance_Message, err
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

// Package smartmoney
package transaction

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/summary"
	"github.com/blue-jay/blueprint/model/transaction"
	"github.com/peddlrph/lib/utilities"

	"github.com/blue-jay/core/router"
)

var (
	uri = "/transaction"
)

// Load the routes.
func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri+"/create", Store, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:id", Edit, c...)
	router.Patch(uri+"/edit/:id", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	sum_items, _, err := summary.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		sum_items = []summary.Item{}
	}
	defaultFormat := "Mon, 02-Jan-2006"

	for i := 0; i < len(sum_items); i++ {
		sum_items[i].SnapshotDate_Formatted = sum_items[i].SnapshotDate.Time.Format(defaultFormat)
		//items[i].Details_Split = strings.Split(items[i].Details, "|")
		sum_items[i].Cash_String = utilities.DisplayPrettyNullFloat64(sum_items[i].Cash)
		sum_items[i].Loads_String = utilities.DisplayPrettyNullFloat64(sum_items[i].Loads)
		sum_items[i].SmartMoney_String = utilities.DisplayPrettyNullFloat64(sum_items[i].SmartMoney)
		sum_items[i].Codes_String = utilities.DisplayPrettyNullFloat64(sum_items[i].Codes)
		sum_items[i].Total_String = utilities.DisplayPrettyNullFloat64(sum_items[i].Total)
	}

	items, sum, _, err1 := transaction.All(c.DB)
	if err1 != nil {
		c.FlashErrorGeneric(err1)
		items = []transaction.Item{}
	}

	//defaultFormat := "Mon, 02-Jan-2006"

	for i := 0; i < len(items); i++ {
		items[i].Trans_Datetime_Formatted = items[i].Trans_Datetime.Time.Format(defaultFormat)
		items[i].Details_Split = strings.Split(items[i].Details, "|")
		items[i].Amount_String = utilities.DisplayPrettyNullFloat64(items[i].Amount)
	}

	v := c.View.New("transaction/index")
	v.Vars["items"] = items
	v.Vars["sum_items"] = sum_items
	v.Vars["sum"] = sum
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("transaction/create")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	fmt.Println(r.FormValue("trans_datetime"))

	_, err := transaction.Create(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))
	//_, err := transaction.Create(c.DB, r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

// Show displays a single item.
func Show(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := transaction.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("transaction/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := transaction.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("transaction/edit")
	c.Repopulate(v.Vars, "amount")
	v.Vars["item"] = item
	v.Vars["setdate"] = item.Trans_Datetime.Time.Format("2006-01-02")
	v.Render(w, r)
}

// Update handles the edit form submission.
func Update(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("amount") {
		Edit(w, r)
		return
	}

	_, err := transaction.Update(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"), c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Edit(w, r)
		return
	}

	c.FlashSuccess("Item updated.")
	c.Redirect(uri)
}

// Destroy handles the delete form submission.
func Destroy(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	_, err := transaction.DeleteSoft(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Package cash
package expenses

import (
	//"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/expenses"

	"github.com/blue-jay/core/router"
	"github.com/peddlrph/lib/utilities"
)

var (
	uri = "/expenses"
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

	items, sum, _, err := expenses.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		items = []expenses.Item{}
	}

	defaultFormat := "Mon, 02-Jan-2006"

	for i := 0; i < len(items); i++ {
		items[i].Trans_Datetime_Formatted = items[i].Trans_Datetime.Time.Format(defaultFormat)
		items[i].Details_Split = strings.Split(items[i].Details, "|")
		items[i].Amount_String = utilities.DisplayPrettyNullFloat64(items[i].Amount)
	}

	v := c.View.New("expenses/index")
	v.Vars["items"] = items
	v.Vars["sum"] = utilities.DisplayPrettyFloat(sum)
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("expenses/create")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "amount")
	v.Render(w, r)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	//_, err1 := strconv.ParseFloat(r.FormValue("amount"), 64)
	//if err1 != nil {
	//	c.FlashErrorGeneric(err1)
	// items = []cash.Item{}
	//}
	//fmt.Println(r.FormValue("amount"), IsNumeric(r.FormValue("amount")))

	//if !c.FormValid("amount") && !IsNumeric(r.FormValue("amount")) {
	if !IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := expenses.Spend(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))

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

	item, _, err := expenses.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("expenses/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := expenses.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("expenses/edit")
	c.Repopulate(v.Vars, "amount")
	v.Vars["item"] = item
	v.Vars["setdate"] = item.Trans_Datetime.Time.Format("2006-01-02")
	//fmt.Println(item.Trans_Datetime.Time.Format("2006-01-02"))
	v.Render(w, r)
}

// Update handles the edit form submission.
func Update(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("amount") {
		Edit(w, r)
		return
	}

	_, err := expenses.Update(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"), c.Param("id"))
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

	_, err := expenses.DeleteSoft(c.DB, c.Param("id"))
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

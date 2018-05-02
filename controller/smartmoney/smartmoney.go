// Package smartmoney
package smartmoney

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/smartmoney"
	"github.com/peddlrph/lib/utilities"

	"github.com/blue-jay/core/router"
)

var (
	uri = "/smartmoney"
)

// Load the routes.
func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri+"/create", Store, c...)
	router.Get(uri+"/smartmoneywithcash", SmartMoneyWithCash, c...)
	router.Post(uri+"/smartmoneywithcash", SmartMoneyWithCashSave, c...)
	router.Get(uri+"/encashsmartmoney", EncashSmartMoney, c...)
	router.Post(uri+"/encashsmartmoney", EncashSmartMoneySave, c...)
	router.Get(uri+"/transfertovirtual", TransferToVirtual, c...)
	router.Post(uri+"/transfertovirtual", TransferToVirtualSave, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:id", Edit, c...)
	router.Patch(uri+"/edit/:id", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	items, sum, _, err := smartmoney.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		items = []smartmoney.Item{}
	}

	defaultFormat := "Mon, 02-Jan-2006"

	for i := 0; i < len(items); i++ {
		items[i].Trans_Datetime_Formatted = items[i].Trans_Datetime.Time.Format(defaultFormat)
		items[i].Details_Split = strings.Split(items[i].Details, "|")
		items[i].Amount_String = utilities.DisplayPrettyNullFloat64(items[i].Amount)
	}

	v := c.View.New("smartmoney/index")
	v.Vars["items"] = items
	v.Vars["sum"] = utilities.DisplayPrettyFloat(sum)
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("smartmoney/create")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func SmartMoneyWithCash(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("smartmoney/smartmoneywithcash")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "amount")
	v.Render(w, r)
}

func SmartMoneyWithCashSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := smartmoney.SmartMoneyWithCash(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func EncashSmartMoney(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("smartmoney/encashsmartmoney")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "amount")
	v.Render(w, r)
}

func EncashSmartMoneySave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := smartmoney.EncashSmartMoney(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func TransferToVirtual(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("smartmoney/transfertovirtual")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "amount")
	v.Render(w, r)
}

func TransferToVirtualSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := smartmoney.TransferToVirtual(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	//if !c.FormValid("name") {
	//	Create(w, r)
	//	return
	//}

	if !IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := smartmoney.Create(c.DB, r.FormValue("amount"), r.FormValue("details"))
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

	item, _, err := smartmoney.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("smartmoney/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := smartmoney.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("smartmoney/edit")
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

	_, err := smartmoney.Update(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"), c.Param("id"))
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

	_, err := smartmoney.DeleteSoft(c.DB, c.Param("id"))
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

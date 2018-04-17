// Package code
package code

import (
	"net/http"
	"strings"
	"time"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/code"

	"github.com/blue-jay/core/router"
	u "github.com/peddlrph/lib/utilities"
)

var (
	uri = "/code"
)

// code the routes.
func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri+"/create", Store, c...)
	router.Get(uri+"/buycode", BuyCode, c...)
	router.Post(uri+"/buycode", BuyCodeSave, c...)
	router.Get(uri+"/code2muni", Code2Muni, c...)
	router.Post(uri+"/code2muni", Code2MuniSave, c...)
	router.Get(uri+"/code2brgy", Code2Brgy, c...)
	router.Post(uri+"/code2brgy", Code2BrgySave, c...)
	router.Get(uri+"/code2dealer", Code2Dealer, c...)
	router.Post(uri+"/code2dealer", Code2DealerSave, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:id", Edit, c...)
	router.Patch(uri+"/edit/:id", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	items, sum, cnt, _, err := code.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		items = []code.Item{}
	}

	defaultFormat := "Mon, 02-Jan-2006"

	for i := 0; i < len(items); i++ {
		items[i].Trans_Datetime_Formatted = items[i].Trans_Datetime.Time.Format(defaultFormat)
		items[i].Details_Split = strings.Split(items[i].Details, "|")
		items[i].Amount_String = u.DisplayPrettyNullFloat64(items[i].Amount)
	}

	prettysum := u.DisplayPrettyFloat(sum)
	prettycnt := u.DisplayPrettyFloat(cnt)

	v := c.View.New("code/index")
	v.Vars["items"] = items
	v.Vars["sum"] = prettysum
	v.Vars["cnt"] = prettycnt
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("code/create")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func BuyCode(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("code/buycode")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func Code2Muni(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("code/code2muni")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func Code2Brgy(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("code/code2brgy")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func Code2Dealer(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("code/code2dealer")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func BuyCodeSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	//code_cost := 3698

	if !u.IsPositiveInteger(r.FormValue("code_count")) {
		c.FlashNotice("Enter valid amount")
		BuyCode(w, r)
		return
	}

	_, err := code.BuyCode(c.DB, r.FormValue("trans_datetime"), r.FormValue("code_count"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		BuyCode(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func Code2MuniSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !u.IsPositiveInteger(r.FormValue("code_count")) {
		c.FlashNotice("Enter valid amount")
		Code2Muni(w, r)
		return
	}

	_, err := code.Code2Muni(c.DB, r.FormValue("trans_datetime"), r.FormValue("code_count"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Code2Muni(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func Code2BrgySave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !u.IsPositiveInteger(r.FormValue("code_count")) {
		c.FlashNotice("Enter valid amount")
		BuyCode(w, r)
		return
	}

	_, err := code.Code2Brgy(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Code2Brgy(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func Code2DealerSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !u.IsPositiveInteger(r.FormValue("code_count")) {
		c.FlashNotice("Enter valid amount")
		Code2Dealer(w, r)
		return
	}

	_, err := code.Code2Dealer(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Code2Dealer(w, r)
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

	if !u.IsPositiveInteger(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := code.Create(c.DB, r.FormValue("amount"), r.FormValue("details"))
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

	item, _, err := code.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("code/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := code.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("code/edit")
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

	_, err := code.Update(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"), c.Param("id"))
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

	_, err := code.DeleteSoft(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}

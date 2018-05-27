// Package load
package load

import (
	//"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/load"
	"github.com/blue-jay/blueprint/model/phonebook"

	"github.com/blue-jay/core/router"
	"github.com/peddlrph/lib/utilities"
)

var (
	uri = "/load"
)

// Load the routes.
func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri+"/create", Store, c...)
	router.Get(uri+"/replenishwithcash", ReplenishWithCash, c...)
	router.Post(uri+"/replenishwithcash", ReplenishWithCashSave, c...)
	router.Get(uri+"/replenishwithsm", ReplenishWithSM, c...)
	router.Post(uri+"/replenishwithsm", ReplenishWithSMSave, c...)
	router.Get(uri+"/load2muni", Load2Muni, c...)
	router.Post(uri+"/load2muni", Load2MuniSave, c...)
	router.Get(uri+"/load2brgy", Load2Brgy, c...)
	router.Post(uri+"/load2brgy", Load2BrgySave, c...)
	router.Get(uri+"/load2dealer", Load2Dealer, c...)
	router.Post(uri+"/load2dealer", Load2DealerSave, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:id", Edit, c...)
	router.Patch(uri+"/edit/:id", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	items, sum, _, err := load.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		items = []load.Item{}
	}

	defaultFormat := "Mon, 02-Jan-2006"

	for i := 0; i < len(items); i++ {
		items[i].Trans_Datetime_Formatted = items[i].Trans_Datetime.Time.Format(defaultFormat)
		items[i].Details_Split = strings.Split(items[i].Details, "|")
		items[i].Amount_String = utilities.DisplayPrettyNullFloat64(items[i].Amount)
	}

	//vv := strconv.FormatFloat(sum, 'f', 2, 32)

	prettysum := utilities.DisplayPrettyFloat(sum)

	//fmt.Printf("%T: %f: %.2f: %v\n", sum, sum, sum, pretty)

	v := c.View.New("load/index")
	v.Vars["items"] = items
	v.Vars["sum"] = prettysum

	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("load/create")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

// Create displays the create form.
func ReplenishWithCash(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("load/replenishwithcash")
	v.Vars["curdate"] = now.Format("2006-01-02")
	//c.Repopulate(v.Vars, "amount")
	v.Render(w, r)
}

func ReplenishWithSM(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("load/replenishwithsm")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func Load2Muni(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("load/load2muni")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func Load2Brgy(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("load/load2brgy")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func Load2Dealer(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	now := time.Now()

	v := c.View.New("load/load2dealer")
	v.Vars["curdate"] = now.Format("2006-01-02")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

func ReplenishWithCashSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := load.ReplenishWithCash(c.DB, r.FormValue("trans_datetime"), r.FormValue("mobile_num"), r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func ReplenishWithSMSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := load.ReplenishWithSM(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func Load2MuniSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) || !utilities.IsMobileNumber(r.FormValue("mobile_number")) {
		c.FlashNotice("Enter valid entries")
		Load2Muni(w, r)
		return
	}

	_, err := load.Load2Muni(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("mobile_number"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Load2Muni(w, r)
		return
	}
	_, _ = phonebook.AddMobileNumberToPhonebook(c.DB, r.FormValue("mobile_number"))

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func Load2BrgySave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) || !utilities.IsMobileNumber(r.FormValue("mobile_number")) {
		c.FlashNotice("Enter valid entries")
		Load2Brgy(w, r)
		return
	}

	_, err := load.Load2Brgy(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("mobile_number"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Load2Brgy(w, r)
		return
	}
	_, _ = phonebook.AddMobileNumberToPhonebook(c.DB, r.FormValue("mobile_number"))

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

func Load2DealerSave(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !utilities.IsNumeric(r.FormValue("amount")) || !utilities.IsMobileNumber(r.FormValue("mobile_number")) {
		c.FlashNotice("Enter valid entries")
		Load2Dealer(w, r)
		return
	}

	_, err := load.Load2Dealer(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("mobile_number"), r.FormValue("details"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Load2Dealer(w, r)
		return
	}
	_, _ = phonebook.AddMobileNumberToPhonebook(c.DB, r.FormValue("mobile_number"))

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

	if !utilities.IsNumeric(r.FormValue("amount")) {
		c.FlashNotice("Enter valid amount")
		Create(w, r)
		return
	}

	_, err := load.Create(c.DB, r.FormValue("amount"), r.FormValue("details"))
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

	item, _, err := load.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("load/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := load.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("load/edit")
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

	_, err := load.Update(c.DB, r.FormValue("trans_datetime"), r.FormValue("amount"), r.FormValue("details"), c.Param("id"))
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

	_, err := load.DeleteSoft(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}

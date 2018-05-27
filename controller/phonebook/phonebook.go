// Package phonebook
package phonebook

import (
	//	"fmt"
	"net/http"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/phonebook"

	"github.com/blue-jay/core/router"
)

var (
	uri = "/phonebook"
)

// Load the routes.
func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri+"/create", Store, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:mobile_number", Edit, c...)
	router.Patch(uri+"/edit/:mobile_number", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	items, _, err := phonebook.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		items = []phonebook.Item{}
	}

	v := c.View.New("phonebook/index")
	v.Vars["items"] = items
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("phonebook/create")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("name") {
		Create(w, r)
		return
	}

	_, err := phonebook.Create(c.DB, r.FormValue("name"))
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

	item, _, err := phonebook.ByMobileNumber(c.DB, c.Param("mobile_number"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("phonebook/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := phonebook.ByMobileNumber(c.DB, c.Param("mobile_number"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	//fmt.Println(item)

	v := c.View.New("phonebook/edit")
	//c.Repopulate(v.Vars, "mobile_number")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Update handles the edit form submission.
func Update(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("name") {
		Edit(w, r)
		return
	}

	_, err := phonebook.Update(c.DB, c.Param("mobile_number"), r.FormValue("name"), r.FormValue("category"), r.FormValue("location"))
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

	_, err := phonebook.DeleteSoft(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}

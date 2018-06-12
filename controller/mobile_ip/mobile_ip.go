// Package mobile_ip
package mobile_ip

import (
	"fmt"
	//	"io/ioutil"
	"net/http"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/message"
	"github.com/blue-jay/blueprint/model/mobile_ip"
	"github.com/blue-jay/core/router"
	"github.com/peddlrph/lib/smsgateway"
)

var (
	uri         = "/mobile_ip"
	message_uri = "/message"
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

	items, _, err := mobile_ip.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		items = []mobile_ip.Item{}
	}

	v := c.View.New("mobile_ip/index")
	v.Vars["items"] = items
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("mobile_ip/create")
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

	_, err := mobile_ip.Create(c.DB, r.FormValue("name"))
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

	item, _, err := mobile_ip.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("mobile_ip/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := mobile_ip.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("mobile_ip/edit")
	c.Repopulate(v.Vars, "name")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Update handles the edit form submission.
func Update(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("ip_address") {
		Edit(w, r)
		return
	}
	if r.FormValue("sync_code") != "a1b2c3d4e5" {
		Edit(w, r)
		return
	}

	if smsgateway.CheckStatus(r.FormValue("ip_address")) != "ready" {
		fmt.Println("Mobile Device: Offline")
	} else {
		fmt.Println("Mobile Device: ready")
		//mesgs, err := smsgateway.GetMessages(r.FormValue("ip_address"), "1000")
		mesgs, err := smsgateway.GetMessages(r.FormValue("ip_address"), "10000")
		//mesgs, err := ioutil.ReadFile("./asset/messages/messages.json")

		if err != nil {
			fmt.Println("Error retrieving messages.")
			fmt.Println("Upload FAILED")
		} else {
			fmt.Println("Messages retrieved")
			_ = message.SyncMessages(c.DB, mesgs)
		}
		//fmt.Println("Messages Uploaded")
	}

	//_, err := mobile_ip.Update(c.DB, r.FormValue("ip_address"), c.Param("id"))
	//if err != nil {
	//	c.FlashErrorGeneric(err)
	//	Edit(w, r)
	//	return
	//}

	//c.FlashSuccess("Item updated.")
	c.Redirect(message_uri)
}

// Destroy handles the delete form submission.
func Destroy(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	_, err := mobile_ip.DeleteSoft(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}

/*
func EditMobileIP(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := mobile_ip.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("mobile_ip/edit")
	//c.Repopulate(v.Vars, "name")
	v.Vars["item"] = item
	v.Render(w, r)
}

func UpdateMobileIP(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("ip_address") {
		Edit(w, r)
		return
	}

	_, err := mobile_ip.Update(c.DB, r.FormValue("ip_address"), c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		Edit(w, r)
		return
	}

	c.FlashSuccess("IP Address updated.")
	c.Redirect(uri)
}
*/

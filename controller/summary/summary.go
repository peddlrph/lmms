// Package summary
package summary

import (
	//"bytes"
	"bufio"
	//"fmt"
	"net/http"
	"os"
	"time"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/summary"

	"github.com/blue-jay/core/router"
	"github.com/peddlrph/lib/utilities"

	"github.com/wcharczuk/go-chart"
)

var (
	uri = "/summary"
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

	items, _, err := summary.All(c.DB)
	if err != nil {
		c.FlashErrorGeneric(err)
		items = []summary.Item{}
	}
	defaultFormat := "Mon, 02-Jan-2006"

	for i := 0; i < len(items); i++ {
		items[i].SnapshotDate_Formatted = items[i].SnapshotDate.Time.Format(defaultFormat)
		//items[i].Details_Split = strings.Split(items[i].Details, "|")
		items[i].Cash_String = utilities.DisplayPrettyNullFloat64(items[i].Cash)
		items[i].Loads_String = utilities.DisplayPrettyNullFloat64(items[i].Loads)
		items[i].SmartMoney_String = utilities.DisplayPrettyNullFloat64(items[i].SmartMoney)
		items[i].Codes_String = utilities.DisplayPrettyNullFloat64(items[i].Codes)
		items[i].Total_String = utilities.DisplayPrettyNullFloat64(items[i].Total)
	}

	daily_earnings, _, _ := summary.DailyEarnings(c.DB)

	for i := 0; i < len(daily_earnings); i++ {
		daily_earnings[i].Trans_Datetime_Formatted = daily_earnings[i].Trans_Datetime.Time.Format(defaultFormat)
		daily_earnings[i].Amount_String = utilities.DisplayPrettyNullFloat64(daily_earnings[i].Amount)
	}

	//prettysum := utilities.DisplayPrettyFloat(sum)

	pie := chart.PieChart{
		Title:  "Summary",
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			{Value: items[0].Cash.Float64, Label: "Cash"},
			{Value: items[0].Loads.Float64, Label: "Loads"},
			{Value: items[0].SmartMoney.Float64, Label: "SmartMoney"},
			{Value: items[0].Codes.Float64, Label: "Codes"},
		},
	}

	outputfile := "./asset/static/outputfile.png"

	//buffer := []byte{}

	f, _ := os.Create(outputfile)

	writer := bufio.NewWriter(f)

	defer f.Close()

	_ = pie.Render(chart.PNG, writer)

	writer.Flush()

	//_ = ioutil.WriteFile(outputfile, buffer, 0644)
	//check(err)

	currentTime := time.Now()

	v := c.View.New("summary/index")
	v.Vars["items"] = items
	v.Vars["today"] = currentTime.Format(defaultFormat)
	//fmt.Println(daily_earnings)
	//v.Vars["buf"] = outputfile
	v.Vars["daily_earnings"] = daily_earnings
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("summary/create")
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

	_, err := summary.Create(c.DB, r.FormValue("name"))
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

	item, _, err := summary.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("summary/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := summary.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("summary/edit")
	c.Repopulate(v.Vars, "name")
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

	_, err := summary.Update(c.DB, r.FormValue("name"), c.Param("id"))
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

	_, err := summary.DeleteSoft(c.DB, c.Param("id"))
	if err != nil {
		c.FlashErrorGeneric(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}

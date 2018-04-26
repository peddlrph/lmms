// Package summary
package summary

import (
	//"bytes"
	"bufio"
	//"fmt"
	"net/http"
	"os"
	"sort"
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

	daily_earnings, _, _ := summary.DailyEarnings(c.DB, "7")

	m := make(map[string]float64)
	n := make(map[string]float64)

	for i := 0; i < len(daily_earnings); i++ {
		daily_earnings[i].Trans_Datetime_Formatted = daily_earnings[i].Trans_Datetime.Time.Format(defaultFormat)
		daily_earnings[i].Amount_String = utilities.DisplayPrettyNullFloat64(daily_earnings[i].Amount)

		m[daily_earnings[i].Trans_Datetime.Time.Format("2006-01-02")] = m[daily_earnings[i].Trans_Datetime.Time.Format("2006-01-02")] + daily_earnings[i].Amount.Float64
		n[daily_earnings[i].Trans_code] = n[daily_earnings[i].Trans_code] + daily_earnings[i].Amount.Float64
	}

	//fmt.Println(m)

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

	pie2 := chart.PieChart{
		Title:  "Summary",
		Width:  512,
		Height: 512,
		Values: []chart.Value{},
	}

	//outputfile := "./asset/static/outputfile.png"

	//buffer := []byte{}

	f, _ := os.Create(outputfile)

	writer := bufio.NewWriter(f)

	defer f.Close()

	_ = pie.Render(chart.PNG, writer)

	writer.Flush()

	sbc := chart.BarChart{
		Title:      "Daily Total Earnings",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: []chart.Value{},
	}

	//vBars := sbc.Bars
	idx := 0
	// fmt.Println("Length of m:", len(m))
	slice := make([]chart.Value, len(m))

	slice1 := make([]summary.Earning, len(m))
	//fmt.Println("Arr:", slice)
	for k, v := range m {

		//		fmt.Println(k, v)
		//fmt.Printf("type: %T\n", sbc.Bars)
		slice[idx].Label = k
		slice[idx].Value = v

		slice1[idx].Trans_Datetime_Formatted = k
		slice1[idx].Amount_String = utilities.DisplayPrettyFloat64(v)

		idx++
	}

	sort.Slice(slice, func(i, j int) bool { return slice[i].Label < slice[j].Label })

	sbc.Bars = slice

	outputfile2 := "./asset/static/outputfile2.png"

	//buffer := []byte{}

	f2, _ := os.Create(outputfile2)

	writer2 := bufio.NewWriter(f2)

	defer f2.Close()

	_ = sbc.Render(chart.PNG, writer2)

	writer2.Flush()

	idx2 := 0

	slice2 := make([]chart.Value, len(n))

	//make([]summary.Earning, len(n))

	slice3 := make([]summary.Earning, len(n))
	//fmt.Println("Arr:", slice)
	for k, v := range n {
		//sbc.Bars[idx].Label = k
		//sbc.Bars[idx].Value = v

		//		fmt.Println(k, v)
		//fmt.Printf("type: %T\n", sbc.Bars)
		slice2[idx2].Label = k
		slice2[idx2].Value = v

		slice3[idx2].Trans_code = k
		slice3[idx2].Amount_String = utilities.DisplayPrettyFloat64(v)

		idx2++
	}

	pie2.Values = slice2

	outputfile3 := "./asset/static/outputfile3.png"

	f3, _ := os.Create(outputfile3)

	writer3 := bufio.NewWriter(f3)

	defer f3.Close()

	_ = pie2.Render(chart.PNG, writer3)

	writer3.Flush()

	//fmt.Println(n)

	//	fmt.Println(daily_earnings)
	//	fmt.Println(slice3)

	currentTime := time.Now()

	v := c.View.New("summary/index")
	v.Vars["items"] = items
	v.Vars["today"] = currentTime.Format(defaultFormat)
	//fmt.Println(daily_earnings)
	//v.Vars["buf"] = outputfile
	//v.Vars["daily_earnings"] = daily_earnings
	v.Vars["daily_earnings"] = slice1
	v.Vars["earnings_by_transcode"] = slice3
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

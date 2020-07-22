package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	reportsgtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/reportsgtm-v1"

	"github.com/go-echarts/go-echarts/charts"
)

var dataPoint []int64

// DemoConfig struct
type DemoConfig struct {
	Domain   string `json:"domain"`
	Property string `json:"property"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

var demoConfig DemoConfig

// QueryDetails for dig request
type QueryDetails struct {
	Location string
	Hostname string
}

// Page for templates
type Page struct {
	Title   string
	Body    string
	Success bool
}

// Response from APIs
type Response struct {
	Title   string
	Body    string
	Success bool
}

func loadPage(title string) (*Page, error) {
	filename := "static/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(body), Success: true}, nil
}

func loadDemoConfig() {
	// Democonfiguration
	data, err := ioutil.ReadFile("conf/conf.json")
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(data, &demoConfig)
	if err != nil {
		fmt.Print(err)
	}

	// Edgegrid
	config, _ := edgegrid.Init("~/.edgerc", "edgednspoc")
	reportsgtm.Init(config)
	optArgs := make(map[string]string)
	optArgs["start"] = demoConfig.Start
	optArgs["end"] = demoConfig.End
	testPropertyTraffic, err := reportsgtm.GetTrafficPerProperty(demoConfig.Domain, demoConfig.Property, optArgs)
	if err == nil {
		for _, v := range testPropertyTraffic.DataRows {
			dataPoint = append(dataPoint, v.Datacenters[0].Requests)
		}
	}
}

func digHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/base.html", "static/dig.html")
	p, _ := loadPage("dig")

	if r.Method != http.MethodPost {
		t.Execute(w, p)
		return
	}

	queryDetails := QueryDetails{
		Location: r.FormValue("location"),
		Hostname: r.FormValue("hostname"),
	}
	fmt.Println(queryDetails)
	response := Response{Title: "dig", Body: queryDetails.Location, Success: true}
	_ = response

	t.Execute(w, response)
}

func lineHandler(w http.ResponseWriter, _ *http.Request) {
	nameItems := []string{}
	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: "DNS Requests over Time"})
	line.AddXAxis(nameItems).
		AddYAxis("DNS Requests", dataPoint,
			charts.LabelTextOpts{Show: true},
			charts.AreaStyleOpts{Opacity: 0.2},
			charts.LineOpts{Smooth: true})
	f, err := os.Create("static/dnshits.html")
	if err != nil {
		fmt.Println(err)
	}
	line.Render(w, f)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate := template.Must(template.ParseFiles("static/index.html"))
	indexTemplate.Execute(w, nil)
}

func main() {
	loadDemoConfig()
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/dig", digHandler)
	http.HandleFunc("/dnshits", lineHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	http.ServeFile(w, r, "static/dig.html")
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
	f, err := os.Create("line.html")
	if err != nil {
		fmt.Println(err)
	}
	line.Render(w, f)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	loadDemoConfig()
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/dig", digHandler)
	http.HandleFunc("/line", lineHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

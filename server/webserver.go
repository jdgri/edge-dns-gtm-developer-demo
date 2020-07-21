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
	data, err := ioutil.ReadFile("conf/conf.json")
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(data, &demoConfig)
	if err != nil {
		fmt.Print(err)
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
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

	config, _ := edgegrid.Init("~/.edgerc", "edgednspoc")
	reportsgtm.Init(config)
	optArgs := make(map[string]string)
	optArgs["start"] = "2020-07-18T00:00:00Z"
	optArgs["end"] = "2020-07-18T23:59:59Z"
	testPropertyTraffic, err := reportsgtm.GetTrafficPerProperty("edgedns.zone.akadns.net", "mirror-failover", optArgs)
	if err == nil {
		for k, v := range testPropertyTraffic.DataRows {
			fmt.Println("Time period: ", k)
			fmt.Println("Time stamp: ", v.Timestamp)
			fmt.Println("Requests: ", v.Datacenters[0].Requests)
			dataPoint = append(dataPoint, v.Datacenters[0].Requests)
		}
	}

	http.HandleFunc("/line", handler)
	http.HandleFunc("/", viewHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

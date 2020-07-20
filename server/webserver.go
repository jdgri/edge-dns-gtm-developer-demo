package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	reportsgtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/reportsgtm-v1"

	"github.com/go-echarts/go-echarts/charts"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	nameItems := []string{"A", "B", "C", "D", "E", "F"}
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "Bar"})
	bar.AddXAxis(nameItems).
		AddYAxis("A", []int{20, 30, 40, 10, 24, 36}).
		AddYAxis("B", []int{35, 14, 25, 68, 44, 23})
	f, err := os.Create("bar.html")
	if err != nil {
		fmt.Println(err)
	}
	bar.Render(w, f)
}

func main() {
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
		}
	}

	http.HandleFunc("/", handler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	"os"
)

type alert struct {
	Labels map[string]string
	Annotations map[string]string
	State string
	Activeat string
	Value string
}

type alertList struct {
	Status string
	Data map[string][]alert
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: missing parameter")
		fmt.Println("Usage: prom-alerts-filter <PrometheusInstance> <AlertName>")
		os.Exit(1)
	}

	var prometheusInstance = os.Args[1]
	var alertName = ""

	if len(os.Args) >= 3 {
		alertName = os.Args[2]
	}

	jsonResponse := getActiveAlerts(prometheusInstance)

	var alerts alertList
	json.Unmarshal([]byte(jsonResponse), &alerts)

	var count int
	if len(alertName) > 0 {
		count = getFilteredAlerts(&alerts, alertName)
	} else {
		count = getAllFiringAlerts(&alerts)
	}

	if count == 0 {
		fmt.Printf("No alerts found in %v\n", prometheusInstance )
	}

}

func getActiveAlerts(prometheusURL string) string {
	URL := strings.Replace("https://URL/api/v1/alerts", "URL", prometheusURL, 1)
	resp, err := http.Get(URL) // TODO: parametrise
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	return string(body)
}

func getFilteredAlerts(alerts *alertList, alertName string) int {
	count := 0
	for k, v := range alerts.Data["alerts"] {
		if v.Labels["alertname"] == alertName {
			json, err := json.MarshalIndent(alerts.Data["alerts"][k], "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(json))
			count++
		}
	}

	return count
}

func getAllFiringAlerts(alerts *alertList) int {
	count := 0
	for k := range alerts.Data["alerts"] {
		fmt.Println(alerts.Data["alerts"][k].Labels["alertname"])
		count++
	}
	return count
}

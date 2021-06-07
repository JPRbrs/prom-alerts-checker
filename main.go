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

	prometheusInstance := os.Args[1]
	alertName := os.Args[2]

	jsonResponse := getActiveAlerts(prometheusInstance)

	var alerts alertList
	json.Unmarshal([]byte(jsonResponse), &alerts)

	getFiringAlerts(&alerts, alertName)
}

func getActiveAlerts(prometheusURL string) string {
	// See first answer here
	// https://stackoverflow.com/questions/17156371/how-to-get-json-response-from-http-get
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

func getFiringAlerts(alerts *alertList, alertName string) {

	// fmt.Println(alerts.Data["alerts"][0].Labels["alertname"])

	for k, v := range alerts.Data["alerts"] {
		// fmt.Printf("key: %v, value: %v,\n", k, v.Labels["alertname"])
		//fmt.Println(v.Labels["alertname"], alertName)
		if v.Labels["alertname"] == alertName {
			json, err := json.MarshalIndent(alerts.Data["alerts"][k], "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(json))
		}

	}
}

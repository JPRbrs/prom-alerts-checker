package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
)

type Alert struct {
	Labels map[string]string
	Annotations map[string]string
	State string
	Activeat string
	Value string
}

type AlertList struct {
	Status string
	Data map[string][]Alert
}


func main() {
	jsonResponse := getActiveAlerts("")

	var alerts AlertList
	json.Unmarshal([]byte(jsonResponse), &alerts)

	processAlerts(&alerts)
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

func processAlerts(alerts *AlertList) {

	// fmt.Println(alerts.Data["alerts"][0].Labels["alertname"])

	for k, v := range alerts.Data["alerts"] {
		fmt.Printf("key: %v, value: %v\n", k, v.Labels["alertname"])
	}
}

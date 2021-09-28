package main

import (
	"errors"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
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
	if len(os.Args) < 3 {
		fmt.Println("Error: missing parameter")
		fmt.Println("Usage: ./prom-alerts-filter <channel> <environment>> [AlertName]")
		fmt.Println("Example ./prom-alerts-filter devops ltitdev")
		os.Exit(1)
	}

	var prometheusInstance = os.Args[1]
	var team = os.Args[1]
	var env = os.Args[2]
	var alertName = ""

	if len(os.Args) == 4 {
		alertName = os.Args[3]
	}

	var prometheusURL = fmt.Sprintf("https://%s-k8s-prometheus.%s.corp-apps.com/api/v1/alerts", team, env)

	jsonResponse, err := getActiveAlerts(prometheusURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var alerts alertList
	json.Unmarshal([]byte(jsonResponse), &alerts)

	mapa := make(map[string]int)
	count := 0
	if len(alertName) > 0 {
		count = getFilteredAlerts(&alerts, alertName)
	} else {
		count, mapa = getAllFiringAlerts(&alerts, mapa)
		fmt.Println(fmt.Sprintf("Total number of alerts for %s: %v", prometheusURL, count))
		for k, v := range mapa {
			fmt.Println(fmt.Sprintf("Alert name: %s, instances: %v", k, v))
		}
	}

	if count == 0 {
		fmt.Printf("No alerts found in %v\n", prometheusInstance )
	}

}

func getActiveAlerts(prometheusURL string) (string, error) {
	resp, err := http.Get(prometheusURL)
	if err != nil {
		return "", errors.New(fmt.Sprintf(
			"Error when calling Prometheus API. \nPlease check URL is correct: %s", prometheusURL))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading Prometheus response")
	}

	return string(body), nil
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

func getAllFiringAlerts(alerts *alertList, mapa map[string]int) (int, map[string]int) {
	count := 0
	for _, v := range alerts.Data["alerts"] {
		// fmt.Println(alerts.Data["alerts"][k].Labels["alertname"])
		// fmt.Println(v.Labels["alertname"])
		_, mapContainsKey := mapa[v.Labels["alertname"]]
		if mapContainsKey {
			mapa[v.Labels["alertname"]]++
		} else {
			mapa[v.Labels["alertname"]] = 1
		}
		count++
	}
	return count, mapa
}

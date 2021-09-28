# Prometheus-alerts-filter

This small Golang tool was born out of the necessity of tracking Prometheus alerts in
Kubernetes non-production Mattermost channel. The amount of alerts in that channels makes
virtually impossible to keep track of alerts when debugging them

# Installation

Clone the repo, navigate into it and execute go build .
You can also install it in your go path with go install .

Go 1.16.3 was used to develope it but it should run in most modern versions. Please report
if your version of Go does not allow to run/install/compile it

# Usage
Usage: prom-alerts-filter <channel> <environment>> [AlertName]
Example prom-alerts-filter devops ltitdev

AlertName is optional, if not passed, it will return a list of the alerts in the given Prometheus

# Collaboration
This script can (and hopefully will) be improved, so please let me know your
ideas/requests or if you want to tinker with Go, feel free to add PRs

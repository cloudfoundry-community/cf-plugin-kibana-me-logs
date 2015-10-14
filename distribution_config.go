package main

import (
	"os"
)

// If you fork and change any of these, remember to create your own releases (see README)

// kibanaMeLogsRepo is the repo to clone when auto-creating kibana-me-logs UI apps
func kibanaMeLogsRepo() string {
	var repo = os.Getenv("KIBANA_ME_LOGS_REPO")
	if repo == "" {
		repo = "https://github.com/cloudfoundry-community/kibana-me-logs"
	}
	return repo
}

package daemon

import "github.com/aks067/devtrack/internal/tracker"

func RunDaemon() {
	tracker.Start(".")
}

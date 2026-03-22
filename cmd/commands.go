package cmds

import (
	"devtrack/internal/github"
	"devtrack/internal/tracker"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func RunDaemon() {
	tracker.Start(".")
}

func RunStatus() {
	var logFile tracker.LogFile
	content, err := os.ReadFile("data/logs.json")
	if err != nil {
		if os.IsNotExist(err) {
			logFile = tracker.LogFile{Sessions: []tracker.Session{}}
		} else {
			log.Fatal("Could not read json file")
		}
	}
	json.Unmarshal(content, &logFile)
	items, err := githubTrack.FetchCommits("aks067")
	if err != nil {
		log.Fatal("Error when fecthing commits")
	}
	maxDuration := int64(0)
	for _, session := range logFile.Sessions {
		if session.Duration > maxDuration {
			maxDuration = session.Duration
		}
	}
	sort.Slice(logFile.Sessions, func(i, j int) bool {
		return logFile.Sessions[i].Duration > logFile.Sessions[j].Duration
	})
	m := githubTrack.CountCommitsByProject(items)
	for _, session := range logFile.Sessions {
		d := time.Duration(session.Duration) * time.Second
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		barWidth := 30
		filled := int(session.Duration * int64(barWidth) / maxDuration)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  📁 %-15s ⏱ %dh%02dm   📦 %d commits\n", session.Project, hours, minutes, m[session.Project])
		fmt.Printf("  %s\n", bar)
	}
}

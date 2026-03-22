package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

type LogFile struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	Project   string    `json:"project"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
	Duration  int64     `json:"duration_sec"`
}

var currentSession *Session

const timeout time.Duration = 10 * time.Minute

var inactivityTimer *time.Timer
var lastEventTime time.Time

func resetTimer() {
	if inactivityTimer == nil {
		inactivityTimer = time.AfterFunc(timeout, func() { closeSession() })
	} else {
		inactivityTimer.Reset(timeout)
	}
}

func getProjectName(projectsDir, filePath string) string {
	path, err := filepath.Rel(projectsDir, filePath)
	if err != nil {
		log.Fatal("Could not retrieve porjects folder")
	}

	projectFolder := strings.Split(path, "/")
	return projectFolder[0]
}

func startSession(project string) {
	var newSession Session
	newSession.Project = project
	newSession.StartTime = time.Now()
	currentSession = &newSession
}

func saveSession(session Session) {
	var logFile LogFile
	content, err := os.ReadFile("data/logs.json")
	if err != nil {
		if os.IsNotExist(err) {
			logFile = LogFile{Sessions: []Session{}}
		} else {
			log.Fatal("Could not read json file")
		}
	}
	json.Unmarshal(content, &logFile)
	found := false
	for i, sess := range logFile.Sessions {
		if sess.Project == session.Project {
			logFile.Sessions[i].StartTime = session.StartTime
			logFile.Sessions[i].EndTime = session.EndTime
			logFile.Sessions[i].Duration += session.Duration
			found = true
		}
	}
	if !found {
		logFile.Sessions = append(logFile.Sessions, session)
	}
	data, err := json.MarshalIndent(logFile, "", " ")
	os.WriteFile("data/logs.json", data, 0644)
}

func closeSession() {
	currentSession.EndTime = time.Now()
	currentSession.Duration = int64(currentSession.EndTime.Sub(currentSession.StartTime).Seconds())
	fmt.Printf("Project %s, saved!\n", currentSession.Project)
	saveSession(*currentSession)
	currentSession = nil
}

func Start(projectsDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v\n", err)
	}
	defer watcher.Close()
	filepath.WalkDir(projectsDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() == true {
			err = watcher.Add(path)
			if err != nil {
				log.Fatalf("Failed to add file to watcher: %v\n", err)
			}
		}
		return nil
	})

	context, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Write) {
				if time.Since(lastEventTime) < 100*time.Millisecond {
					continue
				}
				lastEventTime = time.Now()
				projectName := getProjectName(projectsDir, event.Name)
				if currentSession == nil {
					startSession(projectName)
				} else if currentSession != nil && currentSession.Project != projectName {
					closeSession()
					startSession(projectName)
				}
				resetTimer()
				fmt.Printf("Project: %s, file change at %s\n", projectName, event.Name)
			}
		case err := <-watcher.Errors:
			log.Printf("Watcher error: %v\n", err)
		case <-context.Done():
			fmt.Println()
			fmt.Println("Program Interupted...")
			if currentSession != nil {
				fmt.Println("Saving Session...")
				closeSession()
			}
			return
		}
	}
}

package tracker

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func getProjectName(projectsDir, filePath string) string {
	path, err := filepath.Rel(projectsDir, filePath)
	if err != nil {
		log.Fatal("Could not retrieve porjects folder")
	}

	projectFolder := strings.Split(path, "/")
	return projectFolder[0]
}

func Start(projectsDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()
	filepath.WalkDir(projectsDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() == true {
			err = watcher.Add(path)
			if err != nil {
				log.Fatalf("Failed to add file to watcher: %v", err)
			}
		}
		return nil
	})

	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Write) {
				log.Printf("File changed in ./%s", event.Name)
			}
		case err := <-watcher.Errors:
			log.Printf("Watcher error: %v", err)
		}
	}
}

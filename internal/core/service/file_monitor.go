package service

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

var (
	hotFolder    = "./hot"
	backupFolder = "./backup"
	logFile      = "demo4_cli_golang.out"

	fs = afero.NewOsFs() // Default to real filesystem, but we can override for tests
)

func SetFileSystem(newFs afero.Fs) {
	fs = newFs
}

func StartMonitoring() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watcher.Add(hotFolder); err != nil {
		log.Fatal(err)
	}

	log.Println("Monitoring started on:", hotFolder)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			processEvent(event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func processEvent(event fsnotify.Event) {
	fileName := filepath.Base(event.Name)

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		log.Println("File created:", event.Name)
		logAction("CREATED", event.Name)
		if strings.HasPrefix(fileName, "delete_") {
			deleteFile(event.Name)
		} else {
			backupFile(event.Name)
		}

	case event.Op&fsnotify.Write == fsnotify.Write:
		log.Println("File modified:", event.Name)
		logAction("MODIFIED", event.Name)
		if strings.HasPrefix(fileName, "delete_") {
			deleteFile(event.Name)
		} else {
			backupFile(event.Name)
		}

	case event.Op&fsnotify.Remove == fsnotify.Remove:
		log.Println("File deleted:", event.Name)
		logAction("DELETED", event.Name)

	case event.Op&fsnotify.Rename == fsnotify.Rename:
		log.Println("File renamed:", event.Name)
		logAction("RENAMED", event.Name)
		renameFile(event.Name)
	}
}

func backupFile(filePath string) {
	dest := filepath.Join(backupFolder, filepath.Base(filePath)+".bak")
	input, err := afero.ReadFile(fs, filePath)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	if err := afero.WriteFile(fs, dest, input, 0644); err != nil {
		log.Println("Error creating backup:", err)
		return
	}

	logAction("BACKUP", filePath)
}

func deleteFile(filePath string) {
	hotPath := filepath.Join(hotFolder, filepath.Base(filePath))
	hotPathWithoutPrefix := strings.Replace(hotPath, "delete_", "", 1)
	backupPath := filepath.Join(backupFolder, filepath.Base(hotPathWithoutPrefix)+".bak")

	if exists, _ := afero.Exists(fs, hotPath); exists {
		if err := fs.Remove(hotPath); err != nil {
			log.Println("Error deleting file with 'delete_' prefix from hot folder:", err)
			// log.error("Error deleting file from hot folder:", err)
		}
	}

	if exists, _ := afero.Exists(fs, hotPathWithoutPrefix); exists {
		if err := fs.Remove(hotPathWithoutPrefix); err != nil {
			log.Println("Error deleting file from hot folder:", err)
		}
	}

	if exists, _ := afero.Exists(fs, backupPath); exists {
		if err := fs.Remove(backupPath); err != nil {
			log.Println("Error deleting file from backup folder:", err)
		}
	}

	logAction("DELETE", filePath)
}

// implement renameFile that will delete the previous file and create a new file with the same content
func renameFile(filePath string) {
	hotPath := filepath.Join(hotFolder, filepath.Base(filePath))
	backupPath := filepath.Join(backupFolder, filepath.Base(hotPath)+".bak")

	// delete the file
	if err := fs.Remove(backupPath); err != nil {
		log.Println("Error deleting file from backup folder:", err)
		return
	}

	logAction("RENAME", filePath)
}

func logAction(action, filePath string) {
	logEntry := time.Now().Format("2006-01-02 15:04:05") + " " + action + " " + filePath + "\n"

	file, err := fs.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error writing log:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(logEntry)
	if err != nil {
		log.Println("Error writing log entry:", err)
	}
}

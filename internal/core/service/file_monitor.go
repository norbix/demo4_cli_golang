package service

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	hotFolder    = "./hot"
	backupFolder = "./backup"
	logFile      = "demo4_cli_golang.out"
)

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

	if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
		if strings.HasPrefix(fileName, "delete_") {
			deleteFile(event.Name)
		} else {
			backupFile(event.Name)
		}
	}

	if event.Op&fsnotify.Rename == fsnotify.Rename {
		log.Println("File renamed:", event.Name)
		if strings.HasPrefix(fileName, "delete_") {
			deleteFile(event.Name)
		}
	}
}

func backupFile(filePath string) {
	dest := filepath.Join(backupFolder, filepath.Base(filePath)+".bak")
	input, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	if err := os.WriteFile(dest, input, 0644); err != nil {
		log.Println("Error creating backup:", err)
		return
	}

	logAction("BACKUP", filePath)
}

func deleteFile(filePath string) {
	hotPath := filepath.Join(hotFolder, filepath.Base(filePath))
	backupPath := filepath.Join(backupFolder, filepath.Base(filePath)+".bak")

	if _, err := os.Stat(hotPath); err == nil {
		if err := os.Remove(hotPath); err != nil {
			log.Println("Error deleting file from hot folder:", err)
		}
	}

	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Remove(backupPath); err != nil {
			log.Println("Error deleting file from backup folder:", err)
		}
	}

	logAction("DELETE", filePath)
}

func logAction(action, filePath string) {
	logEntry := time.Now().Format("2006-01-02 15:04:05") + " " + action + " " + filePath + "\n"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error writing log:", err)
		return
	}
	defer file.Close()

	file.WriteString(logEntry)
}

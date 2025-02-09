package service

import (
	"encoding/json"
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
	stateFile    = "demo4_cli_golang.state.json"

	fs = afero.NewOsFs()
)

type AppState struct {
	Files map[string]time.Time `json:"files"`
}

func SetFileSystem(newFs afero.Fs) {
	fs = newFs
}

func saveState(state AppState) {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Println("Error marshalling state:", err)
		return
	}

	if err := afero.WriteFile(fs, stateFile, data, 0644); err != nil {
		log.Println("Error saving state:", err)
	}
}

func loadState() AppState {
	state := AppState{Files: make(map[string]time.Time)}
	if exists, _ := afero.Exists(fs, stateFile); !exists {
		return state
	}

	data, err := afero.ReadFile(fs, stateFile)
	if err != nil {
		log.Println("Error reading state file:", err)
		return state
	}

	if err := json.Unmarshal(data, &state); err != nil {
		log.Println("Error unmarshalling state:", err)
	}

	return state
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

	state := loadState()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			processEvent(event, &state)
			saveState(state)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func processEvent(event fsnotify.Event, state *AppState) {
	fileName := filepath.Base(event.Name)
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		log.Println("File created:", event.Name)
		logAction("CREATED", event.Name)
		// state.Files[event.Name] = time.Now()
		if strings.HasPrefix(fileName, "delete_") {
			deleteFile(event.Name)
		} else {
			backupFile(event.Name)
		}
		state.Files[event.Name] = time.Now()

	case event.Op&fsnotify.Write == fsnotify.Write:
		log.Println("File modified:", event.Name)
		logAction("MODIFIED", event.Name)
		// state.Files[event.Name] = time.Now()
		if strings.HasPrefix(fileName, "delete_") {
			deleteFile(event.Name)
		} else {
			backupFile(event.Name)
		}
		state.Files[event.Name] = time.Now()

	case event.Op&fsnotify.Remove == fsnotify.Remove:
		log.Println("File deleted:", event.Name)
		logAction("DELETED", event.Name)
		delete(state.Files, event.Name)

	case event.Op&fsnotify.Rename == fsnotify.Rename:
		log.Println("File renamed:", event.Name)
		logAction("RENAMED", event.Name)
		renameFile(event.Name, state)
		delete(state.Files, event.Name)
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

func renameFile(filePath string, state *AppState) {
	hotPath := filepath.Join(hotFolder, filepath.Base(filePath))
	backupPath := filepath.Join(backupFolder, filepath.Base(hotPath)+".bak")

	if err := fs.Remove(backupPath); err != nil {
		log.Println("Error deleting file from backup folder:", err)
	}

	state.Files[filePath] = time.Now()
	saveState(*state)

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

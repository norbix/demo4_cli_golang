package service

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setupTestFs() afero.Fs {
	testFs := afero.NewMemMapFs() // Virtual filesystem
	SetFileSystem(testFs)

	// Create hot and backup directories
	_ = testFs.MkdirAll(hotFolder, 0755)
	_ = testFs.MkdirAll(backupFolder, 0755)

	return testFs
}

func TestBackupFile(t *testing.T) {
	testFs := setupTestFs()
	filePath := filepath.Join(hotFolder, "testfile.txt")

	// Create a dummy file in hot folder
	_ = afero.WriteFile(testFs, filePath, []byte("test content"), 0644)

	backupFile(filePath)

	backupPath := filepath.Join(backupFolder, "testfile.txt.bak")
	exists, _ := afero.Exists(testFs, backupPath)

	assert.True(t, exists, "Backup file should exist")
}

func TestDeleteFile(t *testing.T) {
	testFs := setupTestFs()
	filePath := filepath.Join(hotFolder, "delete_testfile.txt")
	backupPath := filepath.Join(backupFolder, "testfile.txt.bak")

	// Create dummy files in both hot and backup folders
	_ = afero.WriteFile(testFs, filePath, []byte("delete this"), 0644)
	_ = afero.WriteFile(testFs, backupPath, []byte("backup content"), 0644)

	deleteFile(filePath)

	hotExists, _ := afero.Exists(testFs, filePath)
	backupExists, _ := afero.Exists(testFs, backupPath)

	assert.False(t, hotExists, "File should be deleted from hot folder")
	assert.False(t, backupExists, "File should be deleted from backup folder")
}

func TestRenameFile(t *testing.T) {
	testFs := afero.NewMemMapFs()
	SetFileSystem(testFs)

	filePath := filepath.Join(hotFolder, "rename_testfile.txt")
	backupPath := filepath.Join(backupFolder, "rename_testfile.txt.bak")

	// Create a dummy backup file
	_ = afero.WriteFile(testFs, backupPath, []byte("backup content"), 0644)

	// Ensure backup file exists before renaming
	backupExistsBefore, _ := afero.Exists(testFs, backupPath)
	assert.True(t, backupExistsBefore, "Backup file should exist before renaming")

	// Initialize a test state
	state := &AppState{Files: make(map[string]time.Time)}

	// Run rename
	renameFile(filePath, state)

	// Check existence after renaming
	backupExistsAfter, _ := afero.Exists(testFs, backupPath)

	// Assertions
	assert.False(t, backupExistsAfter, "Backup file should be deleted after renaming")
}

func TestLoadState(t *testing.T) {
	testFs := setupTestFs()                               // Ensures consistent file system
	fs = testFs                                           // Make sure `fs` in loadState() uses the test filesystem
	stateFile = filepath.Join(backupFolder, "state.json") // Ensure path matches `loadState()`

	// Case 1: State file does not exist
	t.Run("State file does not exist", func(t *testing.T) {
		state := loadState()
		assert.Empty(t, state.Files, "Expected an empty AppState when state file is missing")
	})

	// Case 2: State file contains valid JSON
	t.Run("State file contains valid JSON", func(t *testing.T) {
		expectedState := AppState{
			Files: map[string]time.Time{
				"file1.txt": time.Now(),
			},
		}

		// Serialize expected state to JSON
		jsonData, _ := json.Marshal(expectedState)

		// Ensure we are writing to the correct file path
		err := afero.WriteFile(fs, stateFile, jsonData, 0644)
		assert.NoError(t, err, "Error writing state file")

		// **Force filesystem sync to avoid buffering issues**
		file, err := fs.Open(stateFile)
		assert.NoError(t, err, "Error opening state file for sync")
		file.Sync()
		file.Close()

		// Ensure file exists before calling `loadState()`
		exists, err := afero.Exists(fs, stateFile)
		assert.NoError(t, err, "Error checking if state file exists")
		assert.True(t, exists, "State file should exist before loading")

		// Load state from file
		state := loadState()

		// **Debugging Output**
		t.Logf("Loaded state: %+v", state)

		// Assert that the state was loaded correctly
		assert.Equal(t, len(expectedState.Files), len(state.Files), "Expected state to be properly populated")
		_, existsInState := state.Files["file1.txt"]
		assert.True(t, existsInState, "Expected 'file1.txt' in loaded state")
	})
}

package service

import (
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

//
// func TestRenameFile(t *testing.T) {
// 	testFs := afero.NewMemMapFs()
// 	SetFileSystem(testFs)
//
// 	filePath := filepath.Join(hotFolder, "rename_testfile.txt")
// 	backupPath := filepath.Join(backupFolder, "rename_testfile.txt.bak")
//
// 	// Create a dummy backup file
// 	_ = afero.WriteFile(testFs, backupPath, []byte("backup content"), 0644)
//
// 	// Ensure backup file exists before renaming
// 	backupExistsBefore, _ := afero.Exists(testFs, backupPath)
// 	assert.True(t, backupExistsBefore, "Backup file should exist before renaming")
//
// 	// Run rename
// 	renameFile(filePath)
//
// 	// Check existence after renaming
// 	backupExistsAfter, _ := afero.Exists(testFs, backupPath)
//
// 	// Assertions
// 	assert.False(t, backupExistsAfter, "Backup file should be deleted after renaming")
// }

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

package service

import (
	"path/filepath"
	"testing"

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

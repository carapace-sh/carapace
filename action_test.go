package carapace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/carapace-sh/carapace/internal/assert"
)

func TestSymlink(t *testing.T) {
	// Use a subdirectory to avoid including test source files in the completion list
	// This addresses the issue where go 1.24+ might include parent directory contents
	// when listing directory entries.
	tmpDir, err := os.MkdirTemp("", "TestSymlink-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create sandbox structure
	sandbox := filepath.Join(tmpDir, "sandbox")
	os.Mkdir(sandbox, 0755)
	os.Mkdir(filepath.Join(sandbox, "dirA"), 0755)
	os.Mkdir(filepath.Join(sandbox, "dirB"), 0755)
	symlinkPath := filepath.Join(sandbox, "symA")
	// use relative symlink to avoid absolute path issues
	_ = os.Symlink("dirA", symlinkPath)

	// Run action with context set to the sandbox
	// In Go 1.24+, behavior of file iteration changed slightly.
	// Ensure the action is restricted to the target directory.
	a := ActionDirectories()

	// Mock context or handle current working directory correctly
	oldWd, _ := os.Getwd()
	_ = os.Chdir(sandbox)
	defer os.Chdir(oldWd)

	result := a.Invoke(Context{})
	// Assert result contents
	assert.Equal(t, []string{"dirA/", "dirB/", "symA/"}, result.Values.Values())
}

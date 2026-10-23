package util

import (
	"os"
	"path/filepath"
	"testing"
)

// canonical resolves symlinks and short names (windows) and normalizes separators so that
// paths built from tempdirs match paths reported by the filesystem (e.g. /var -> /private/var on macOS).
func canonical(t *testing.T, path string) string {
	t.Helper()
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	return filepath.ToSlash(filepath.Clean(path))
}

func TestFindReverseSymlinkedDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "carapace-util-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repoDir := filepath.Join(tmpDir, "myrepo")
	subDir := filepath.Join(repoDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(repoDir, ".git"), 0755); err != nil {
		t.Fatalf("failed to create .git: %v", err)
	}

	symlinkDir := filepath.Join(tmpDir, "symlink-test")
	if err := os.Symlink(subDir, symlinkDir); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	found, err := FindReverse(symlinkDir, ".git")
	if err != nil {
		t.Fatalf("FindReverse failed from symlinked dir: %v", err)
	}

	expected := canonical(t, filepath.Join(repoDir, ".git"))
	if canonical(t, found) != expected {
		t.Errorf("expected %q, got %q", expected, found)
	}
}

func TestFindReverseNormalDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "carapace-util-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repoDir := filepath.Join(tmpDir, "myrepo")
	subDir := filepath.Join(repoDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(repoDir, ".git"), 0755); err != nil {
		t.Fatalf("failed to create .git: %v", err)
	}

	found, err := FindReverse(subDir, ".git")
	if err != nil {
		t.Fatalf("FindReverse failed from normal dir: %v", err)
	}

	expected := canonical(t, filepath.Join(repoDir, ".git"))
	if canonical(t, found) != expected {
		t.Errorf("expected %q, got %q", expected, found)
	}
}

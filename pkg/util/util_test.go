package util

import (
	"os"
	"path/filepath"
	"testing"
)

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

	expected := filepath.Join(repoDir, ".git")
	if found != expected {
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

	expected := filepath.Join(repoDir, ".git")
	if found != expected {
		t.Errorf("expected %q, got %q", expected, found)
	}
}

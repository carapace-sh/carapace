package traverse

import (
	"os"
	"path/filepath"
	"testing"
)

type mockContext struct {
	dir string
	env []string
}

func (m mockContext) Abs(s string) (string, error) {
	if s == "" {
		if m.dir == "" {
			return filepath.Abs("")
		}
		return filepath.Abs(m.dir)
	}
	return filepath.Abs(filepath.Join(m.dir, s))
}

func (m mockContext) Getenv(key string) string {
	for _, e := range m.env {
		if len(e) > len(key)+1 && e[:len(key)+1] == key+"=" {
			return e[len(key)+1:]
		}
	}
	return ""
}

func (m mockContext) LookupEnv(key string) (string, bool) {
	v := m.Getenv(key)
	return v, v != ""
}

func TestParentSymlinkedDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "carapace-traverse-test-*")
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

	ctx := mockContext{dir: symlinkDir, env: os.Environ()}

	result, err := Parent(".git")(ctx)
	if err != nil {
		t.Fatalf("Parent(.git) failed from symlinked dir: %v", err)
	}

	expected := filepath.ToSlash(repoDir)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestParentNormalDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "carapace-traverse-test-*")
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

	ctx := mockContext{dir: subDir, env: os.Environ()}

	result, err := Parent(".git")(ctx)
	if err != nil {
		t.Fatalf("Parent(.git) failed from normal dir: %v", err)
	}

	expected := filepath.ToSlash(repoDir)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestGitWorkTreeSymlinkedDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "carapace-traverse-test-*")
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

	ctx := mockContext{dir: symlinkDir, env: os.Environ()}

	result, err := GitWorkTree(ctx)
	if err != nil {
		t.Fatalf("GitWorkTree failed from symlinked dir: %v", err)
	}

	expected := filepath.ToSlash(repoDir)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestGitDirSymlinkedDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "carapace-traverse-test-*")
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

	ctx := mockContext{dir: symlinkDir, env: os.Environ()}

	result, err := GitDir(ctx)
	if err != nil {
		t.Fatalf("GitDir failed from symlinked dir: %v", err)
	}

	expected := filepath.ToSlash(repoDir + "/.git")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

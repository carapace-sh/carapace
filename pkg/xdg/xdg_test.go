package xdg

import (
	"runtime"
	"testing"
)

func TestDataDirs(t *testing.T) {
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skip("unix-specific test")
	case "darwin":
		t.Setenv("XDG_DATA_DIRS", "")
		dirs, err := DataDirs()
		if err != nil {
			t.Fatal(err)
		}
		if len(dirs) != 1 || dirs[0] != "/Library/Application Support" {
			t.Fatalf("unexpected dirs: %v", dirs)
		}
	default:
		t.Setenv("XDG_DATA_DIRS", "/data1:relative:/data2")
		dirs, err := DataDirs()
		if err != nil {
			t.Fatal(err)
		}
		if len(dirs) != 2 || dirs[0] != "/data1" || dirs[1] != "/data2" {
			t.Fatalf("expected absolute entries only, got: %v", dirs)
		}

		t.Setenv("XDG_DATA_DIRS", "")
		dirs, err = DataDirs()
		if err != nil {
			t.Fatal(err)
		}
		if len(dirs) != 2 || dirs[0] != "/usr/local/share" || dirs[1] != "/usr/share" {
			t.Fatalf("unexpected defaults: %v", dirs)
		}
	}
}

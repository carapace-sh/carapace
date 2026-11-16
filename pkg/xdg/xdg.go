package xdg

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// UserCacheDir returns the user cache base directory.
func UserCacheDir() (dir string, err error) {
	if dir = os.Getenv("XDG_CACHE_HOME"); dir == "" || !filepath.IsAbs(dir) {
		dir, err = os.UserCacheDir()
	}
	dir = filepath.ToSlash(dir)
	return
}

// UserConfigDir returns the user config base directory.
func UserConfigDir() (dir string, err error) {
	if dir = os.Getenv("XDG_CONFIG_HOME"); dir == "" || !filepath.IsAbs(dir) {
		dir, err = os.UserConfigDir()
	}
	dir = filepath.ToSlash(dir)
	return
}

// ConfigDirs returns the global config base directories.
func ConfigDirs() (dirs []string, err error) {
	switch runtime.GOOS {
	case "windows":
		dir, ok := os.LookupEnv("PROGRAMDATA")
		if !ok {
			return nil, errors.New("missing PROGRAMDATA environment variable")
		}
		dirs = append(dirs, dir)

	case "darwin":
		dirs = append(dirs, "/Library/Application Support")
	default:
		dirs = append(dirs, "/etc/xdg")

	}

	if v, ok := os.LookupEnv("XDG_CONFIG_DIRS"); ok {
		xdgDirs := make([]string, 0)
		for dir := range strings.SplitSeq(v, string(os.PathSeparator)) {
			if filepath.IsAbs(dir) {
				xdgDirs = append(xdgDirs, dir)
			}
		}
		dirs = append(xdgDirs, dirs...)
	}

	for index, dir := range dirs {
		dirs[index] = filepath.ToSlash(dir)
	}
	return
}

// DataDirs returns the system data base directories, ordered by preference.
// On unix this respects $XDG_DATA_DIRS and falls back to the XDG default.
func DataDirs() (dirs []string, err error) {
	switch runtime.GOOS {
	case "windows":
		dir, ok := os.LookupEnv("PROGRAMDATA")
		if !ok {
			return nil, errors.New("missing PROGRAMDATA environment variable")
		}
		dirs = append(dirs, dir)

	case "darwin":
		dirs = append(dirs, "/Library/Application Support")

	default:
		if v, ok := os.LookupEnv("XDG_DATA_DIRS"); ok {
			for dir := range strings.SplitSeq(v, string(os.PathListSeparator)) {
				if filepath.IsAbs(dir) {
					dirs = append(dirs, dir)
				}
			}
		}
		if len(dirs) == 0 {
			dirs = append(dirs, "/usr/local/share", "/usr/share")
		}
	}

	for index, dir := range dirs {
		dirs[index] = filepath.ToSlash(dir)
	}
	return
}

// CarapacePathEntry is a resolved top-level carapace root directory with its group.
type CarapacePathEntry struct {
	Path  string
	Group string // "user" or "system"
}

// CarapacePaths returns ordered carapace top-level directories using XDG:
// user = UserConfigDir()/carapace (group "user"), then system roots from
// DataDirs() with "carapace" appended (group "system"), in order.
func CarapacePaths() ([]CarapacePathEntry, error) {
	userBase, err := UserConfigDir()
	if err != nil {
		return nil, err
	}
	userBase = filepath.ToSlash(userBase)

	var out []CarapacePathEntry
	out = append(out, CarapacePathEntry{Path: filepath.ToSlash(filepath.Join(filepath.FromSlash(userBase), "carapace")), Group: "user"})

	dataDirs, err := DataDirs()
	if err != nil {
		return nil, err
	}
	for _, d := range dataDirs {
		out = append(out, CarapacePathEntry{Path: filepath.ToSlash(filepath.Join(filepath.FromSlash(d), "carapace")), Group: "system"})
	}
	return out, nil
}

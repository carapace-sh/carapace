package xdg

import "os"

// CacheDir returns the cache base directory for carapace.
func CacheDir() (dir string, err error) {
	if dir = os.Getenv("XDG_CACHE_HOME"); dir == "" {
		dir, err = os.UserCacheDir()
	}
	if err == nil {
		dir += "/carapace"
	}
	return
}

// ConfigDir returns the config base directory for carapace.
func ConfigDir() (dir string, err error) {
	if dir = os.Getenv("XDG_CONFIG_HOME"); dir == "" {
		dir, err = os.UserConfigDir()
	}
	if err == nil {
		dir += "/carapace"
	}
	return
}

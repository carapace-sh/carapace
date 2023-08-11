package traverse

import (
	"path/filepath"
)

// GitDir returns the location of the .git folder.
func GitDir(tc Context) (string, error) {
	if dir, ok := tc.LookupEnv("GIT_DIR"); ok {
		return filepath.ToSlash(dir), nil
	}
	return Parent(".git")(tc)
}

// GitWorkTree returns the location of the root of the working directory for a non-bare repository.
func GitWorkTree(tc Context) (string, error) {
	if dir, ok := tc.LookupEnv("GIT_WORK_TREE"); ok {
		return filepath.ToSlash(dir), nil
	}

	gitDir, err := GitDir(tc)
	return filepath.Dir(gitDir), err
}

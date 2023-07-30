package traverse

import (
	"path/filepath"
)

// GitDir is the location of the .git folder.
func GitDir(tc Context) (string, error) {
	if dir, ok := tc.LookupEnv("GIT_DIR"); ok {
		return filepath.ToSlash(dir), nil
	}
	return Parent(".git")(tc)
}

// GitWorkTree is the location of the root of the working directory for a non-bare repository.
func GitWorkTree(tc Context) (string, error) {
	if dir, ok := tc.LookupEnv("GIT_WORK_TREE"); ok {
		return filepath.ToSlash(dir), nil
	}

	if _, ok := tc.LookupEnv("GIT_DIR"); ok {
		wd, _ := tc.Abs("") // TODO provide better access to Context.Dir
		return wd, nil
	}

	gitdir, err := GitDir(tc)
	return filepath.Dir(gitdir), err
}

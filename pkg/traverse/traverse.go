package traverse

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Context interface {
	Abs(s string) (string, error)
	Getenv(key string) string
	LookupEnv(key string) (string, bool)
}

func Parent(names ...string) func(tc Context) (string, error) {
	return func(tc Context) (string, error) {
		wd, err := tc.Abs("")
		if err != nil {
			return "", err
		}

		for _, name := range names {
			if dir, err := traverse(wd, name); err == nil {
				return dir, nil
			}
		}
		formattedNames := fmt.Sprintf("%#v", names)
		formattedNames = strings.TrimPrefix(formattedNames, "[]string{")
		formattedNames = strings.TrimSuffix(formattedNames, "}")
		return "", errors.New("could not find parent directory containing any of: " + formattedNames)
	}
}

// TODO also stop at `~`
func traverse(path string, name string) (target string, err error) {
	var absPath string
	if absPath, err = filepath.Abs(path); err == nil {
		target = filepath.ToSlash(absPath + "/" + name)
		if _, err = os.Stat(target); err != nil {
			parent := filepath.Dir(absPath)
			if parent != path {
				return traverse(parent, name)
			} else {
				err = errors.New("could not find: " + name)
			}
		}
	}
	return
}

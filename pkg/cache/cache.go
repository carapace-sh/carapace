package cache

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"strings"
)

// Key provides a cache key
type Key func() (string, error)

// String creates a CacheKey for given strings
func String(s ...string) Key {
	return func() (string, error) {
		return strings.Join(s, "\n"), nil
	}
}

// FileChecksum creates a CacheKey for given file
func FileChecksum(file string) Key {
	return func() (checksum string, err error) {
		var content []byte
		if content, err = ioutil.ReadFile(file); err == nil {
			checksum = fmt.Sprintf("%x", sha1.Sum(content))
		}
		return
	}
}

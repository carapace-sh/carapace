package cache

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"strings"
)

type CacheKey func() (string, error)

func String(s ...string) CacheKey {
	return func() (string, error) {
		return strings.Join(s, "\n"), nil
	}
}

func FileChecksum(file string) CacheKey {
	return func() (checksum string, err error) {
		var content []byte
		if content, err = ioutil.ReadFile(file); err == nil {
			checksum = fmt.Sprintf("%x", sha1.Sum(content))
		}
		return
	}
}

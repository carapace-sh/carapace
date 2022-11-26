// Package cache provides disk cache for Actions
package cache

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/rsteube/carapace/pkg/cache"
)

// Write persistests given values to file as json.
func Write(file string, rawValues []common.RawValue) (err error) {
	var m []byte
	if m, err = json.Marshal(rawValues); err == nil {
		err = os.WriteFile(file, m, 0600)
	}
	return
}

// Load loads values from file unless modification date exceeds timeout.
func Load(file string, timeout time.Duration) (rawValues []common.RawValue, err error) {
	var content []byte
	var stat os.FileInfo
	if stat, err = os.Stat(file); os.IsNotExist(err) || (timeout > 0 && stat.ModTime().Add(timeout).Before(time.Now())) {
		err = errors.New("not exists or timeout exceeded")
	} else {
		if content, err = os.ReadFile(file); err == nil {
			err = json.Unmarshal(content, &rawValues)
		}
	}
	return
}

// CacheDir creates a cache folder for current user and returns the path.
func CacheDir(name string) (dir string, err error) {
	var userCacheDir string
	if userCacheDir, err = os.UserCacheDir(); err == nil {
		dir = fmt.Sprintf("%v/carapace/%v/%v", userCacheDir, uid.Executable(), name)
		err = os.MkdirAll(dir, 0700)
	}
	return
}

// File returns the cache filename for given values
// TODO cleanup
func File(callerFile string, callerLine int, keys ...cache.Key) (file string, err error) {
	uid := uidKeys(callerFile, strconv.Itoa(callerLine))
	ids := make([]string, 0)
	for _, key := range keys {
		id, err := key()
		if err != nil {
			return "", err
		}
		ids = append(ids, id)
	}
	if dir, err := CacheDir(uid); err == nil {
		file = dir + "/" + uidKeys(ids...)
	}
	return
}

func uidKeys(keys ...string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(keys, "\001"))))
}

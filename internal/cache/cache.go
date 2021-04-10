package cache

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/rsteube/carapace/pkg/cache"
)

// Write persistests given values to file as json
func Write(file string, rawValues []common.RawValue) (err error) {
	var m []byte
	if m, err = json.Marshal(rawValues); err == nil {
		err = ioutil.WriteFile(file, m, 0600)
	}
	return
}

// Load loads values from file unless modification date exceeds timeout
func Load(file string, timeout time.Duration) (rawValues []common.RawValue, err error) {
	var content []byte
	var stat os.FileInfo
	if stat, err = os.Stat(file); os.IsNotExist(err) || (timeout > 0 && stat.ModTime().Add(timeout).Before(time.Now())) {
		err = errors.New("not exists or timeout exceeded")
	} else {
		if content, err = ioutil.ReadFile(file); err == nil {
			err = json.Unmarshal(content, &rawValues)
		}
	}
	return
}

// TempDir creates a temporary folder for current user and returns the path
func TempDir(name string) (dir string, err error) {
	var u *user.User
	if u, err = user.Current(); err == nil {
		dir = fmt.Sprintf("%v/carapace", os.TempDir())
		if err = os.MkdirAll(dir, 0777); err == nil {
			dir = fmt.Sprintf("%v/%v/%v/%v", dir, u.Username, uid.Executable(), name)
			err = os.MkdirAll(dir, 0700)
		}
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
	if dir, err := TempDir(uid); err == nil {
		file = dir + "/" + uidKeys(ids...)
	}
	return
}

func uidKeys(keys ...string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(keys, "\001"))))
}

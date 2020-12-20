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

func Write(file string, rawValues []common.RawValue) (err error) {
	var m []byte
	if m, err = json.Marshal(rawValues); err == nil {
		err = ioutil.WriteFile(file, m, 0600)
	}
	return
}

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

// TODO cleanup
func File(callerFile string, callerLine int, keys ...cache.CacheKey) (file string, err error) {
	uid := Uid(callerFile, strconv.Itoa(callerLine))
	ids := make([]string, 0)
	for _, key := range keys {
		if id, err := key(); err != nil {
			return "", err
		} else {
			ids = append(ids, id)
		}
	}
	if dir, err := TempDir(uid); err == nil {
		file = dir + "/" + Uid(ids...)
	}
	return
}

func Uid(keys ...string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(keys, "\001"))))
}

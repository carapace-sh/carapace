package cache

import (
	"runtime"
	"time"

	"github.com/rsteube/carapace/internal/cache"
	"github.com/rsteube/carapace/pkg/cache/key"
)

// Cache caches a function for given duration and keys.
func Cache(timeout time.Duration, keys ...key.Key) func(f func() ([]byte, error)) ([]byte, error) {
	return func(f func() ([]byte, error)) ([]byte, error) {
		_, file, line, _ := runtime.Caller(1)
		cacheFile, err := cache.File(file, line, keys...)
		if err != nil {
			return nil, err
		}

		content, err := cache.Load(cacheFile, timeout)
		if err != nil {
			content, err = f()
			if err != nil {
				return nil, err
			}
			return content, cache.Write(cacheFile, content)
		}
		return content, nil
	}
}

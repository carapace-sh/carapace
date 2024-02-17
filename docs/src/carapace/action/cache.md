# Cache

[`Cache`] caches an [Action] for a given duration.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	return carapace.ActionValues(
		time.Now().Format("15:04:05"),
	)
}).Cache(5 * time.Second)
```

![](./cache.cast)

> Caches are implicitly assigned a unique key using [`runtime.Caller`] which can change between releases.


## Key

Additional keys like [`key.String`] can be passed as well.

```go
carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
	switch len(c.Parts) {
	case 0:
		return carapace.ActionValues("one", "two").Suffix("/")
	case 1:
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionValues(
				time.Now().Format("15:04:05"),
			)
		}).Cache(10*time.Second, key.String(c.Parts[0]))
	default:
		return carapace.ActionValues()
	}
})
```

![](./cache-key.cast)


## Location

Cache is written as `json` to [`os.UserCacheDir`] using the [Export] format.

```handlebars
{{cacheDir}}/carapace/{{binary}}/{{callerChecksum}}/{{cacheChecksum}}
```

| ID             | x                                | example                                    |
| ----           | ---                              | ---                                        |
| cacheDir       | os.UserCacheDir                  | `~/.cache/`                                |
| binary         | binary name                      | `carapace`                                 |
| callerChecksum | sha1sum using [`runtime.Caller`] | `89be88b670885d3d7855c7169ad7cfd2816a6c37` |
| cacheChecksum  | sh1sum of given [`CacheKeys`]    | `041858daaaa8b084122d4604a3223315c39edc3e` |

[Action]:../action.md
[`Cache`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.Cache
[`key.String`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/key#String
[`CacheKeys`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/cache#CacheKey
[callback actions]:./defaultActions/actionCallback.md
[Export]:../export.md
[`os.UserCacheDir`]:https://pkg.go.dev/os#UserCacheDir
[`runtime.Caller`]:https://pkg.go.dev/runtime#Caller

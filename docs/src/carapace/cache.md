# Cache

[`Cache`] provides a simple way to cache [callback actions].
For this an [InvokedAction] is persisted as `json` to [`os.UserCacheDir`] (see [Export]):

```handlebars
{{cacheDir}}/carapace/{{username}}/{{binary}}/{{callerChecksum}}/{{cacheChecksum}}
```

| ID             | x                                | example                                    |
| ----           | ---                              | ---                                        |
| cacheDir       | os.UserCacheDir                  | `~/.cache/`                                |
| username       | current user                     | `root`                                     |
| binary         | binary name                      | `carapace`                                 |
| callerChecksum | sha1sum using [`runtime.Caller`] | `89be88b670885d3d7855c7169ad7cfd2816a6c37` |
| cacheChecksum  | sh1sum of given [`CacheKeys`]    | `041858daaaa8b084122d4604a3223315c39edc3e` |

[`Cache`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.Cache
[`CacheKeys`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/cache#CacheKey
[callback actions]:./defaultActions/actionCallback.md
[Export]:./export.md
[InvokedAction]:./invokedAction.md
[`os.UserCacheDir`]:https://pkg.go.dev/os#UserCacheDir
[`runtime.Caller`]:https://pkg.go.dev/runtime#Caller




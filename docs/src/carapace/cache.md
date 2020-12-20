# Cache

> Cache is still _experimental_ and will undergo some changes

[`Cache`](https://pkg.go.dev/github.com/rsteube/carapace#Action.Cache) provides a simple way to cache [callback actions](./action/actionCallback.md).
For this the values of an [InvokedAction](./invokedAction.md) are persisted as `json` to [`os.TempDir`](https://pkg.go.dev/os#TempDir):

```handlebars
{{TempDir}}/carapace/{{username}}/{{binary}}/{{callerChecksum}}/{{cacheChecksum}}
```

| ID | x | example |
|----|---|---|
| TempDir | os.TempDir | `/tmp` |
| username | current user | `root` |
| binary | binary name | `carapace` |
| callerChecksum | sha1sum using [`runtime.Caller`](https://pkg.go.dev/runtime#Caller) | `89be88b670885d3d7855c7169ad7cfd2816a6c37` |
| cacheChecksum | sh1sum of given [`CacheKeys`](https://pkg.go.dev/github.com/rsteube/carapace/pkg/cache#CacheKey) | `041858daaaa8b084122d4604a3223315c39edc3e` |


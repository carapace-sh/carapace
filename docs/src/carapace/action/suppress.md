# Suppress

[`Suppress`] suppresses specific error messages using regular expressions.

```go
carapace.Batch(
	docker.ActionContainers(),
	docker.ActionServices().Supress("This node is not a swarm manager"),
	docker.ActionNetworks(),
	docker.ActionVolumes(),
).ToA()
```

[`Suppress`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Suppress

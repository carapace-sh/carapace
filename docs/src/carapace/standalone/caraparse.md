# caraparse

[caraparse] is a helper tool that uses regex to parse gnu help pages. Due to strong inconsistencies between these the results may differ but generally give a good head start.

```sh
docker node update --help | caraparse -n update -p node -s "Update a node"
```


[caraparse]:https://github.com/rsteube/carapace-bin/tree/master/cmd/caraparse

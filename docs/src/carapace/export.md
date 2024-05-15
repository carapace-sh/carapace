# Export

[`Export`] provides a `json` representation of an [InvokedAction].
It is used to exchange completions between commands with [ActionImport] as well as for [Cache].

```go	
type Export struct {
	version  string   `json:"version"`
	messages []string `json:"messages"`
	nospace  string   `json:"nospace"`
	usage    string   `json:"usage"`
	values   []struct {
		value       string `json:"value"`
		display     string `json:"display"`
		description string `json:"description,omitempty"`
		style       string `json:"style,omitempty"`
		tag         string `json:"tag,omitempty"`
	} `json:"values"`
}
```

| Key            | Description                                                    |
|----------------|----------------------------------------------------------------|
| version        | version of `carapace` being used                               | 
| messages       | list of error messages                                         | 
| nospace        | character suffixes that prevent space suffix (`*` matches all) | 
| usage          | usage message                                                  | 
| values         | list of completion values                                      | 
| -              |                                                                | 
|	value          | value to insert                                                |
|	display        | value to display during completion                             |
|	description    | description of the value                                       |
|	style          | style of the value                                             |
|	tag            | tag of the value                                               |

## Example

```sh
example _carapace export example m<TAB>
```

```json
{
  "version": "unknown",
  "messages": [],
  "nospace": "",
  "usage": "",
  "values": [
    {
      "value": "modifier",
      "display": "modifier",
      "description": "modifier example",
      "style": "yellow",
      "tag": "modifier commands"
    },
    {
      "value": "multiparts",
      "display": "multiparts",
      "description": "multiparts example",
      "tag": "other commands"
    }
  ]
}
```

![](./export.cast)


[ActionImport]:./defaultActions/actionImport.md
[Cache]:./action/cache.md
[`Export`]:https://pkg.go.dev/github.com/carapace-sh/carapace/internal/export#Export
[InvokedAction]:./invokedAction.md

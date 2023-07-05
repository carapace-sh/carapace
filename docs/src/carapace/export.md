# Export

[`Export`] provides a `json` representation of an [InvokedAction].
It is used to exchange completions between commands with [ActionImport] as well as for [Cache].

```go	
type Export struct {
	Version  string   `json:"version"`
	Messages []string `json:"messages"`
	Nospace  string   `json:"nospace"`
	Usage    string   `json:"usage"`
	Values   []struct {
		Value       string `json:"value"`
		Display     string `json:"display"`
		Description string `json:"description,omitempty"`
		Style       string `json:"style,omitempty"`
		Tag         string `json:"tag,omitempty"`
	} `json:"values"`
}
```

| Key            | Description                                                    |
|----------------|----------------------------------------------------------------|
| Version        | version of `carapace` being used                               | 
| Messages       | list of error messages                                         | 
| Nospace        | character suffixes that prevent space suffix (`*` matches all) | 
| Usage          | usage message                                                  | 
| Values         | list of completion values                                      | 
| -              |                                                                | 
|	Value          | value to insert                                                |
|	Display        | value to display during completion                             |
|	Description    | description of the value                                       |
|	Style          | style of the value                                             |
|	Tag            | tag of the value                                               |

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
[`Export`]:https://pkg.go.dev/github.com/rsteube/carapace/internal/export#Export
[InvokedAction]:./invokedAction.md

# Export

Export generates json from the command structure and [Action] (for reuse with [ActionImport]).

## Command structure
```sh
example _carapace export
```

```json
{
  "Name": "example",
  "Short": "example completion",
  "Commands": [
    {
      "Name": "action",
      "Short": "action example",
      "Aliases": [
        "alias"
      ],
      "LocalFlags": [
        {
          "Longhand": "count",
          "Shorthand": "c",
          "Usage": "count flag",
          "Type": "count",
          "NoOptDefVal": "+1"
        },
        {
          "Longhand": "directories",
          "Usage": "files flag",
          "Type": "string"
        },
...
```

## Action

```sh
example _carapace export example action --usergroup root:
```

```json
{
  "Version": "v0.14.0",
  "Nospace": true,
  "RawValues": [
    {
      "Value": "root:root",
      "Display": "root",
      "Description": "0"
    },
    {
      "Value": "root:adm",
      "Display": "adm",
      "Description": "999"
    },
    {
      "Value": "root:wheel",
      "Display": "wheel",
      "Description": "998"
    },
...
```

[Action]:./action.md
[ActionImport]:./action/actionImport.md

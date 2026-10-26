Set Shell "fish"
Source "../fish.tape"
Type "example "
Wait+Screen /(?m)example/
Sleep 1
Tab@1 3
Sleep 2
Wait+Screen /subcommand\ \ \ \ \(subcommand\ example\)/

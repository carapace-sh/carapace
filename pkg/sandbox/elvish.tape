Set Theme "Snazzy"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0
# Set TypingSpeed 0

Hide
Set Shell "elvish"
Type "eval (example _carapace|slurp)" Enter
Type "clear" Enter
Show

Type "example "
Sleep 1 Tab
Sleep 300ms
Screenshot out/elvish.png

Sleep 2
Output     out/elvish.gif
Output     out/elvish.ascii

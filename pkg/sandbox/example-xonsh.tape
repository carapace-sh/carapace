Set Theme "Snazzy"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0
# Set TypingSpeed 0

Hide
Set Shell "xonsh"
Type "exec($(example _carapace xonsh))" Enter
Type "clear" Enter
Show

Type "example "
Sleep 1 Tab@1 2
Sleep 300ms
Screenshot out/example-xonsh.png

Sleep 2
Output     out/example-xonsh.gif
Output     out/example-xonsh.ascii

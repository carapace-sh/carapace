package fish

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`function _%v_quote_suffix
  if not commandline -cp | xargs echo 2>/dev/null >/dev/null
    if commandline -cp | sed 's/$/"/'| xargs echo 2>/dev/null >/dev/null
      echo '"'
    else if commandline -cp | sed "s/\$/'/"| xargs echo 2>/dev/null >/dev/null
      echo "'"
    end
  else 
    echo ""
  end
end

function _%v_callback
  commandline -cp | sed "s/\$/"(_%v_quote_suffix)"/" | sed "s/ \$/ ''/" | xargs %v _carapace fish _
end

complete -c %v -f
complete -c '%v' -f -a '(_%v_callback)' -r
`, cmd.Name(), cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}

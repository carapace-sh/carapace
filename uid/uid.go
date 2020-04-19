package uid

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Command(cmd *cobra.Command) string {
	names := make([]string, 0)
	current := cmd
	for {
		names = append(names, current.Name())
		current = current.Parent()
		if current == nil {
			break
		}
	}

	reverse := make([]string, len(names))
	for i, entry := range names {
		reverse[len(names)-i-1] = entry
	}

	return "_" + strings.Join(reverse, "__")
}

func Flag(cmd *cobra.Command, flag *pflag.Flag) string {
	c := cmd
	for c.HasParent() {
		if c.LocalFlags().Lookup(flag.Name) != nil {
			break
		}
		c = c.Parent()
	}
	// TODO ensure flag acually belongs to command (force error)
	// TODO handle unknown flag error
	return fmt.Sprintf("%v##%v", Command(c), flag.Name)
}

func Positional(cmd *cobra.Command, position int) string {
	// TODO complete function
	return fmt.Sprintf("%v#%v", Command(cmd), position)
}

func find(cmd *cobra.Command, uid string) *cobra.Command {
	var splitted []string
	if splitted = strings.Split(uid[1:], "#"); len(splitted) == 0 { // TODO check for empty uid string
		return nil
	}
	c, _, err := cmd.Root().Find(strings.Split(splitted[0], "__")[1:]) // TODO root if jut one arg
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func Executable() string {
	if executable, err := os.Executable(); err != nil {
		return "echo" // safe fallback that should never happen
	} else if filepath.Base(executable) == "cmd.test" {
		return "example" // for `go test -v ./...`
	} else {
		return filepath.Base(executable)
	}
}

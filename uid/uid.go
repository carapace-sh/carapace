package uid

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func Value(cmd *cobra.Command, args []string, uid string) string {
	// TODO assumes cmd is correct
	if strings.Contains(uid, "##") {
		split := strings.Split(uid, "##")
		if flag := cmd.Flag(split[len(split)-1]); flag != nil {
			if flag.Value.Type() == "stringSlice" {
				slice, _ := cmd.Flags().GetStringSlice(split[len(split)-1])
				return strings.Join(slice, ",")
			} else if flag.Value.Type() == "stringArray" {
				slice, _ := cmd.Flags().GetStringArray(split[len(split)-1])
				return strings.Join(slice, ",")
			} else {
				return flag.Value.String()
			}
		}

	} else if strings.Contains(uid, "#") {
		split := strings.Split(uid, "#")
		if index, err := strconv.Atoi(split[len(split)-1]); err == nil {
			if index > 0 {
				index = index - 1
				if len(args)-1 >= index && index >= 0 {
					return args[index]
				}
			} else if len(args) > 0 && os.Args[len(os.Args)-1] != "" {
				return args[len(args)-1]
			}
		}
	}
	return ""
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

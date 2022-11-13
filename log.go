package carapace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/rsteube/carapace/pkg/ps"
)

type _logger struct {
	*log.Logger
}

func (l _logger) PrintArgs(args []string) {
	if m, err := json.Marshal(args); err == nil {
		l.Println(string(m))
	}
}

var logger = _logger{log.New(ioutil.Discard, "", log.Flags())}

func init() {
	if _, enabled := os.LookupEnv("CARAPACE_LOG"); !enabled {
		return
	}

	tmpdir := fmt.Sprintf("%v/carapace", os.TempDir())
	if err := os.MkdirAll(tmpdir, os.ModePerm); err != nil {
		log.Fatal(err.Error())
	}

	file := fmt.Sprintf("%v/%v.log", tmpdir, uid.Executable())
	if logfileWriter, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		log.Fatal(err.Error())
	} else {
		logger = _logger{log.New(logfileWriter, ps.DetermineShell()+" ", log.Flags()|log.Lmsgprefix)}
	}
}

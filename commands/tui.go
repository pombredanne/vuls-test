/* Vuls - Vulnerability Scanner
Copyright (C) 2016  Future Architect, Inc. Japan.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	c "github.com/sadayuki-matsuno/vuls-test/config"
	"github.com/sadayuki-matsuno/vuls-test/report"
	"github.com/google/subcommands"
	"golang.org/x/net/context"
)

// TuiCmd is Subcommand of host discovery mode
type TuiCmd struct {
	lang     string
	debugSQL bool
	dbpath   string
}

// Name return subcommand name
func (*TuiCmd) Name() string { return "tui" }

// Synopsis return synopsis
func (*TuiCmd) Synopsis() string { return "Run Tui view to anayze vulnerabilites" }

// Usage return usage
func (*TuiCmd) Usage() string {
	return `tui:
	tui [-dbpath=/path/to/vuls.sqlite3]

`
}

// SetFlags set flag
func (p *TuiCmd) SetFlags(f *flag.FlagSet) {
	//  f.StringVar(&p.lang, "lang", "en", "[en|ja]")
	f.BoolVar(&p.debugSQL, "debug-sql", false, "debug SQL")

	wd, _ := os.Getwd()

	defaultDBPath := filepath.Join(wd, "vuls.sqlite3")
	f.StringVar(&p.dbpath, "dbpath", defaultDBPath,
		fmt.Sprintf("/path/to/sqlite3 (default: %s)", defaultDBPath))
}

// Execute execute
func (p *TuiCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c.Conf.Lang = "en"
	c.Conf.DebugSQL = p.debugSQL
	c.Conf.DBPath = p.dbpath

	historyID := ""
	if 0 < len(f.Args()) {
		if _, err := strconv.Atoi(f.Args()[0]); err != nil {
			log.Errorf("First Argument have to be scan_histores record ID: %s", err)
			return subcommands.ExitFailure
		}
		historyID = f.Args()[0]
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Errorf("Failed to read stdin: %s", err)
				return subcommands.ExitFailure
			}
			fields := strings.Fields(string(bytes))
			if 0 < len(fields) {
				historyID = fields[0]
			}
		}
	}
	return report.RunTui(historyID)
}

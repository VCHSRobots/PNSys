// --------------------------------------------------------------------
// cmd_show_log.go -- shows the log files
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/log"
	"epic/lib/util"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const gTopic_show_log = `
The show-log command shows (i.e., dumps) the log files on the console. 
The format of the command is:

    show-log line0=nnn maxlines=nnn date=yyyy-mm-dd gin=t/f

All the parameters are optional.  

The line0 parameter controls where in the log file that output starts.
If negative, lines are counted from the end.  If omitted, line0 defaults to -10.  

The maxlines parameter contols the maximumn number of lines output.
If omitted, output is limited to 50 lines.

The date parameters controls while log file to show, and if omitted, the
current day's log file is shown.

The gin parameter controls weither or not to include gin messages.
By default they are not included.

`

const gTopic_show_log_line = `
The show-log-line commands shows one line of a log file.  This command
is useful to see the full contents of a log line when it is truncated on
in a list output.  The format of the command is:

    show-log-line nnn date=yyyy-mm-dd

where nnn is the line number given in the tabular output of show-log.

`

func init() {
	RegistorCmd("show-log", "", "Shows (dumps) log files (see topic).", handle_show_log)
	RegistorCmd("show-log-line", "", "Shows one log message form a log file (see topic).", handle_show_log_line)
	RegistorTopic("show-log", gTopic_show_log)
	RegistorTopic("show-log-line", gTopic_show_log_line)
}

func handle_show_log(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	sline0, _ := util.MapAlias(params, "line0", "0", "start")
	smax, _ := util.MapAlias(params, "maxlines", "max")
	sdate, _ := util.MapAlias(params, "date")
	sgin, _ := util.MapAlias(params, "gin")

	iline0 := -10
	imax := 50
	date := time.Now()
	gin := false
	if !util.Blank(sline0) {
		iline0, err = strconv.Atoi(sline0)
		if err != nil {
			c.Printf("Bad input (%q) for iline0.\n", sline0)
			return
		}
	}
	if !util.Blank(smax) {
		imax, err := strconv.Atoi(smax)
		if err != nil {
			c.Printf("Bad input (%q) for maxlines.\n", smax)
			return
		}
		if imax <= 0 {
			c.Printf("No lines to output. maxlines <= 0. (%s)\n", smax)
			return
		}
	}
	if !util.Blank(sdate) {
		date, err = util.ParseGenericTime(sdate)
		if err != nil {
			c.Printf("Bad input (%s) for date.\n", sdate)
			return
		}
	}
	if !util.Blank(sgin) {
		gin, err = util.StrToBool(sgin, false)
		if err != nil {
			c.Printf("Bad input (%s) for gin.\n", sgin)
			return
		}
	}

	// read entire file.
	content, err := log.ReadLogFile(date)
	if err != nil {
		c.Printf("Error on read: %v\n", err)
		return
	}
	rawlines := strings.Split(content, "\n")
	n := len(rawlines)
	type numline struct {
		line_num int
		text     string
	}
	lines := make([]*numline, 0, len(rawlines))
	for i, x := range rawlines {
		x = strings.TrimSpace(x)
		lines = append(lines, &numline{i, x})
	}
	c.Printf("%d lines in the log file.\n", n)
	if !gin {
		// filter the gin lines.
		tmplns := make([]*numline, 0, len(lines))
		for _, ln := range lines {
			s := strings.TrimSpace(ln.text)
			if !strings.HasPrefix(s, "![GIN]") {
				tmplns = append(tmplns, ln)
			}
		}
		lines = tmplns
		n = len(lines)
		c.Printf("%d lines left after gin messages filtered out.\n", n)
	}
	if iline0 >= n {
		c.Printf("line0 (%d) starts after the end of the file. Nothing to output.\n", iline0)
		return
	}
	if iline0 < 0 {
		iline0 = n + iline0
		if iline0 < 0 {
			iline0 = 0
		}
	}
	c.Printf("Log file of %s, starting at line %d:\n", date.Format("2006-01-02"), iline0)
	tbl := util.NewTable("Line #", "Message")
	nc := 0
	for i := iline0; i < n; i++ {
		snum := fmt.Sprintf("%05d", lines[i].line_num)
		ln := util.FixStrLen(lines[i].text, 125, "...")
		//ln = strings.TrimSpace(ln)
		if len(lines[i].text) > 122 {
			snum += "-"
		} else {
			snum += " "
		}
		tbl.AddRow(snum, ln)
		nc++
		if nc > imax {
			break
		}
	}
	c.Printf("\n%s\n", tbl.Text())
}

func handle_show_log_line(c *Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	if len(args) < 2 {
		c.Printf("Line number not provided.\n")
		return
	}
	iline, err := strconv.Atoi(args[1])
	if err != nil {
		c.Printf("Unrecognizable line number.\n")
		return
	}
	if iline < 0 {
		c.Printf("Negative line numbers don't make sense.\n")
		return
	}

	sdate, _ := util.MapAlias(params, "date")

	date := time.Now()
	if !util.Blank(sdate) {
		date, err = util.ParseGenericTime(sdate)
		if err != nil {
			c.Printf("Bad input (%s) for date.\n", sdate)
			return
		}
	}

	// read entire file.
	content, err := log.ReadLogFile(date)
	if err != nil {
		c.Printf("Error on read: %v\n", err)
		return
	}
	rawlines := strings.Split(content, "\n")
	if iline < 0 || iline >= len(rawlines) {
		c.Printf("Line number out of range. Only %d lines in the file.", len(rawlines))
		return
	}
	txt := rawlines[iline]
	c.Printf("Line %d from log %s:\n%s\n", iline, date.Format("2006-01-02"), txt)
}

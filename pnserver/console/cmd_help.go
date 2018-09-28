// --------------------------------------------------------------------
// cmd_help.go -- Commands for help.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/lib/util"
	"fmt"
	"sort"
	"strings"
)

var gArgs map[string]string
var gTopics map[string]string

func init() {
	RegistorCmd("help", "[topic]", "Gives this help or help on a topic.", handle_help)
	RegistorCmd("?", "", "Condensed help.", handle_help_condensed)
	RegistorArg("topic", "The topic for help. Use 'help topics' to get a list of topics.")
}

func RegistorArg(arg, desc string) {
	if gArgs == nil {
		gArgs = make(map[string]string, 100)
	}
	gArgs[arg] = desc
}

func RegistorTopic(topic, desc string) {
	if gTopics == nil {
		gTopics = make(map[string]string, 100)
	}
	gTopics[topic] = desc
}

func SayHelp() {
	fmt.Printf("%s/n", gHelpText)
}

func show_topic_list() {
	tlst := make([]string, 0, len(gTopics))
	for k, _ := range gTopics {
		tlst = append(tlst, k)
	}
	sort.Strings(tlst)
	sout := "Topic list: "
	for i, t := range tlst {
		if i != 0 {
			sout += ", "
		}
		if len(sout)+len(t) > 80 {
			fmt.Printf("%s\n", sout)
			sout = ""
		}
		sout += t
	}
	if len(sout) > 0 {
		fmt.Printf("%s\n", sout)
	}
}

func handle_help(cmdline string) {
	cwrds := strings.Split(cmdline, " ")
	if len(cwrds) > 1 {
		top := strings.TrimSpace(cwrds[1])
		if top == "topics" {
			show_topic_list()
			return
		}
		desc, ok := gTopics[top]
		if !ok {
			fmt.Printf("Unknown topic. ")
			show_topic_list()
			return
		}
		fmt.Printf("%s\n%s\n", top, desc)
		return
	}
	fmt.Printf("You are now in an interactive command loop. The commnds are:\n\n")
	w1 := 10
	for _, c := range gCmds {
		w1 = max(w1, len(c.CmdName+" "+c.ArgLine+" "))
	}
	w1 = max(w1, len("exit | quit  "))
	for _, c := range gCmds {
		s1 := util.FixStrLen(c.CmdName+" "+c.ArgLine+" ", w1, "")
		fmt.Printf("    %s -- %s\n", s1, c.HelpLine)
	}
	fmt.Printf("    %s -- Exits this program.\n", util.FixStrLen("exit | quit  ", w1, ""))
	fmt.Printf("\n")

	w2 := 10
	argnames := make([]string, 0, len(gArgs))
	for aname, _ := range gArgs {
		w2 = max(w2, len(aname))
		argnames = append(argnames, aname)
	}
	sort.Strings(argnames)
	fmt.Printf("Where the argments are:\n")
	for _, aname := range argnames {
		desc := gArgs[aname]
		sn := util.FixStrLen(aname, w2, "")
		fmt.Printf("    %s -- %s\n", sn, desc)
	}
	fmt.Printf("\n")
}

func handle_help_condensed(cmdline string) {
	sout := ""
	w := 0
	i := 0
	for _, c := range gCmds {
		if c.CmdName == "?" {
			continue
		}
		if i != 0 {
			sout += ", "
			w += 2
		}
		if len(c.CmdName)+w > 80 {
			sout += "\n"
			w = 0
		}
		sout += c.CmdName
		w += len(c.CmdName)
		i += 1
	}
	sout += ", exit"
	fmt.Printf("%s\n", sout)
}

var gHelpText string = `
    pnsserver -- Part Number Server for Epic Robotz -- Sept 2018
    This server runs on a Droplet in the clowd, and assignes part numbers 
    for CAD drawings.
`

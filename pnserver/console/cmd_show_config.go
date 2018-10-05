// --------------------------------------------------------------------
// cmd_show-config.go -- shows the current configuration
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/config"
	"io/ioutil"
	"sort"
)

const gTopic_show_config = `
The show-config command shows the current configuration in memory and also
shows the content of the config.txt file.  The format of the command is:

    show-config

`

func init() {
	RegistorCmd("show-config", "", "Shows the configuration (see topic).", handle_show_config)
	RegistorTopic("show-config", gTopic_show_config)
}

func handle_show_config(c *Context, cmdline string) {
	data, err := ioutil.ReadFile(config.ConfigFileName)
	if err != nil {
		c.Printf("Unable to read config.txt?  Err=%v", err)
	} else {
		c.Printf("Dump of %s:\n", config.ConfigFileName)
		c.Printf("------------------------------------------------ (file start)\n")
		c.Printf("%s\n", string(data))
		c.Printf("------------------------------------------------ (file end)\n")
	}

	m, err := config.GetConfig()
	if err != nil {
		c.Printf("\nUnable to get memory config. Err=%v\n", err)
		return
	}

	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	c.Printf("\nIn Memory Configuration Parameters:\n")
	for _, k := range keys {
		c.Printf("%s=%s\n", k, m[k])
	}
	c.Printf("\n")
}

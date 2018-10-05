// --------------------------------------------------------------------
// cmd_save-config.go -- saves the config file.
//
// Created 2018-10-04 DLB
// --------------------------------------------------------------------

package console

import (
	"epic/pnserver/config"
)

const gTopic_save_config = `
The save-config command saves the current configuration to disk so that on the
next restart of the server the current configuration is used again. This
perserves the following modes and parameters:

    developer bypass mode               -- see dev-bypass 
    universal password mode             -- see universal-pw 
    show log messages on the console    -- see show-log and hide-log
    show gin messages on the console    -- see show-gin and high-gin
    the host address                    -- Not changeable by console command
    the mysql password                  -- Not changeable by console command

The format of the command is:

    save-config

`

func init() {
	RegistorCmd("save-config", "", "Saves the configuration to disk (see topic).", handle_save_config)
	RegistorTopic("save-config", gTopic_save_config)
}

func handle_save_config(c *Context, cmdline string) {
	err := config.WriteConfig()
	if err != nil {
		c.Printf("Err: %v\n", err)
	}
	c.Printf("Success.\n")
}

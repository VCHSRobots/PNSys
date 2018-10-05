// --------------------------------------------------------------------
// pnserver.go -- Main server file for Epics the Part Number System
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

package main

import (
	"epic/lib/log"
	"epic/lib/util"
	"epic/pnserver/config"
	"epic/pnserver/console"
	"epic/pnserver/pages"
	"epic/pnserver/pnsql"
	"epic/pnserver/sessions"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

var gVersion = "Fall 2018 v0.3"
var gServer *gin.Engine
var gHostAddr string = ":8081"

func main() {
	log.Infof("Part Number Server Staring Up. Version: %s", gVersion)
	CheckDirs()
	_, err := config.GetConfig()
	if err != nil {
		err = fmt.Errorf("Unable to get config. Err=%v", err)
		log.Errorf("%v", err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	var ok bool
	gHostAddr, ok = config.GetParam("hostaddr")
	if !ok {
		log.Warnf("The hostaddr config parameter not found. Using %q.\n", gHostAddr)
		gHostAddr = ":8081"
		config.SetParam("hostaddr", gHostAddr)
	}
	sql_pw, ok := config.GetParam("sql_pw")
	if !ok {
		err = fmt.Errorf("Mysql password not found in config file.")
		log.Errorf("%v", err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	err = pnsql.OpenDatabase(sql_pw)
	if err != nil {
		log.Errorf("%v", err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	console.RegistorCmd("version", "", "Gives the version of this server.", handle_version)

	gServer = gin.Default()
	gServer.Static("/css", "./static/css")
	gServer.Static("/img", "./static/img")
	gServer.Static("/js", "./static/js")

	// Load known pages
	plst := pages.GetAllPages()
	for _, p := range plst {
		if p.Invoke == pages.Invoke_GET {
			gServer.GET(p.Route, p.Handlers...)
		}
		if p.Invoke == pages.Invoke_POST {
			gServer.POST(p.Route, p.Handlers...)
		}
	}

	// Configure startup page.
	gServer.GET("/", func(c *gin.Context) { c.Redirect(300, "/NewEpicPN") })

	dev_bypass, ok := config.GetParam("dev_bypass")
	dev_bypass = strings.TrimSpace(dev_bypass)
	if ok && !util.Blank(dev_bypass) {
		fmt.Printf("Developer bypass mode is ON.  For %s as admin.\n", dev_bypass)
		sessions.SetDeveloperBypass(dev_bypass)
	}

	log_on_console, _ := config.GetBoolParam("log_on_console", true)
	if log_on_console {
		fmt.Printf("Showing Log on console. Use 'hide-log' to turn this off.\n")
		log.UseConsole(true)
	}
	config.SetBoolParam("log_on_console", log_on_console)

	gin_on_console, _ := config.GetBoolParam("gin_on_console", true)
	log.AllowPassOnConsole(gin_on_console)
	if !gin_on_console {
		fmt.Printf("GIN messages will NOT be sent to the console termainal.")
	}
	config.SetBoolParam("gin_on_console", gin_on_console)

	allow_universal_pw, _ := config.GetBoolParam("allow_universal_pw", true)
	sessions.SetAllowUniversalPasswords(allow_universal_pw)
	log.Infof("Universal Password mode is %t", allow_universal_pw)
	config.SetBoolParam("allow_universal_pw", allow_universal_pw)

	go RunServer() // Start up and run server in different thread
	fmt.Printf("Server running.  Should be able to access at %s\n", gHostAddr)
	go console.ConsoleLoop() // Process console commands
	<-make(chan int)         // Wait forever here
}

func RunServer() {
	log.Infof("Running Server.  Serving: %s", gHostAddr)
	gServer.Run(gHostAddr)
}

func CheckDirs() {
	paths := []string{"./static", "./static/css", "./static/templates", "./static/img", "./static/js"}
	for _, p := range paths {
		if !util.DirExists(p) {
			err := fmt.Errorf("Static directory (%s) does not exist.", p)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			log.Fatalf("%v", err)
		}
	}
}

func handle_version(c *console.Context, cmdline string) {
	c.Printf("Version: %s\n", gVersion)
}

// --------------------------------------------------------------------
// pnserver.go -- Main server file for Epics the Part Number System
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

// package main

// import (
// 	"net/http"
// 	"strings"
// )

// func sayHello(w http.ResponseWriter, r *http.Request) {
// 	message := r.URL.Path
// 	message = strings.TrimPrefix(message, "/")
// 	message = "Hello " + message
// 	w.Write([]byte(message))
// }
// func main() {
// 	http.HandleFunc("/", sayHello)
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		panic(err)
// 	}
// }

package main

import (
	"epic/lib/log"
	"epic/lib/util"
	"epic/lib/uuid"
	"epic/pnserver/console"
	"epic/pnserver/pages"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var gVersion = "Fall 2018 v0.1"
var gServer *gin.Engine
var gHostAddr string = ":8081"

func main() {
	log.Infof("Part Number Server Staring Up. Version: %s", gVersion)
	CheckDirs()
	console.RegistorCmd("version", "", "Gives the version of this server.", handle_version)

	gServer = gin.Default()
	gServer.Static("/css", "./static/css")
	gServer.Static("/img", "./static/img")
	gServer.Static("/html", "./static/html")
	gServer.GET("/PutCookie", PutCookie)
	gServer.GET("/GetCookie", GetCookie)
	gServer.GET("/Login", pages.Login)
	gServer.POST("/LoginPost", pages.LoginPost)

	// Load known pages
	plst := pages.GetAllPages()
	for _, p := range plst {
		if p.Invoke == pages.Invoke_GET {
			gServer.GET(p.Route, p.Handler)
		}
		if p.Invoke == pages.Invoke_POST {
			gServer.POST(p.Route, p.Handler)
		}
	}

	fmt.Printf("Showing Log to console. Use 'hide-log' to turn this off.\n")
	log.UseConsole(true)

	go RunServer() // Start up and run server in different thread
	fmt.Printf("Server running.  Should be able to access at %s\n", gHostAddr)
	go console.ConsoleLoop() // Process console commands
	<-make(chan int)         // Wait forever here
}

func PutCookie(c *gin.Context) {
	id := uuid.New()
	c.SetCookie("Cred", id.String(), 0, "/", "", false, true)
	c.String(200, "Cookie set: %s for %s", id, c.ClientIP())
}

func GetCookie(c *gin.Context) {
	cookie, err := c.Cookie("Cred")
	if err != nil {
		c.String(200, "Err getting cookie. %v", err)
		return
	}
	c.String(200, "Getting Cookie... Found=%s", cookie)
}

func RunServer() {
	log.Infof("Running Server.  Serving: %s", gHostAddr)
	gServer.Run(gHostAddr)
}

func CheckDirs() {
	paths := []string{"./static", "./static/css", "./static/html", "./static/img",
		"./static/templates"}
	for _, p := range paths {
		if !util.DirExists(p) {
			err := fmt.Errorf("Static directory (%s) does not exist.", p)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			log.Fatalf("%v", err)
		}
	}
}

func handle_version(cmdline string) {
	fmt.Printf("Version: %s\n", gVersion)
}

// func ConsoleLoop() {
// 	gConsole = liner.NewLiner()
// 	defer gConsole.Close()
// 	gConsole.SetCtrlCAborts(false)
// 	if !liner.TerminalSupported() {
// 		fmt.Printf("Terminal not supported. Editting commands won't work.\n")
// 	}
// 	fmt.Printf("Server running.  Should be able to access at %s\n", gHostAddr)
// 	fmt.Printf("Use 'help' for a list of commands.\n")
// 	for {
// 		cmdline, err := gConsole.Prompt("PnSrv> ")
// 		if err == liner.ErrPromptAborted {
// 			continue
// 		}
// 		if err != nil {
// 			fmt.Printf("Input error: %v\n", err)
// 			continue
// 		}
// 		ExecuteCommand(cmdline)
// 		if !util.Blank(cmdline) {
// 			gConsole.AppendHistory(cmdline)
// 		}
// 	}
// }

// func ExecuteCommand(cmdline string) {
// 	cmd := strings.ToLower(strings.TrimSpace(cmdline))
// 	if util.Blank(cmd) {
// 		return
// 	}
// 	switch cmd {
// 	case "quit":
// 		gConsole.Close()
// 		os.Exit(0)
// 		return
// 	case "exit":
// 		gConsole.Close()
// 		os.Exit(0)
// 		return
// 	case "show-log":
// 		log.UseConsole(true)
// 		return
// 	case "hide-log":
// 		log.UseConsole(false)
// 		return
// 	case "help":
// 		Cmd_Help(cmdline)
// 		return
// 	default:
// 		fmt.Printf("Unknown command.\n")
// 		return
// 	}
// }

// func Cmd_Help(cmdline string) {
// 	fmt.Printf("Commands: \n")
// 	fmt.Printf("  help -- List commands\n")
// 	fmt.Printf("  exit -- Shuts down server and exits program.\n")
// }

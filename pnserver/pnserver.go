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
	"epic/pnserver/pnsql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
)

var gVersion = "Fall 2018 v0.3"
var gServer *gin.Engine
var gHostAddr string = ":8081"

func main() {
	log.Infof("Part Number Server Staring Up. Version: %s", gVersion)
	CheckDirs()
	cparams, err := GetConfig("config.txt")
	if err != nil {
		err = fmt.Errorf("Cannot read config.txt file.  %v", err)
		log.Errorf("%v", err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	var ok bool
	gHostAddr, ok = cparams["hostaddr"]
	if !ok {
		log.Warnf("The hostaddr config parameter not found. Using ':8081'.\n")
		gHostAddr = ":8081"
	}
	pw, ok := cparams["pw"]
	if !ok {
		err = fmt.Errorf("Mysql password not found in config file.")
		log.Errorf("%v", err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	err = pnsql.OpenDatabase(pw)
	if err != nil {
		log.Errorf("%v", err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	console.RegistorCmd("version", "", "Gives the version of this server.", handle_version)

	gServer = gin.Default()
	gServer.Static("/css", "./static/css")
	gServer.Static("/img", "./static/img")
	gServer.GET("/PutCookie", PutCookie)
	gServer.GET("/GetCookie", GetCookie)

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
	paths := []string{"./static", "./static/css", "./static/templates", "./static/img"}
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

func GetConfig(filename string) (map[string]string, error) {
	params := make(map[string]string, 10)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return params, err
	}
	lines := strings.Split(string(data), "\n")
	ilinenum := 0
	for _, ln := range lines {
		ilinenum++
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "//") {
			continue
		}
		if util.Blank(ln) {
			continue
		}
		wrds := strings.Split(ln, "=")
		if len(wrds) != 2 {
			return params, fmt.Errorf("Bad syntax on line %d. One equal char not found.\n", ilinenum)
		}
		key := strings.TrimSpace(wrds[0])
		val := strings.TrimSpace(wrds[1])
		params[key] = val
	}
	return params, nil
}

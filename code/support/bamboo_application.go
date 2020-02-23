package support

import (
	"flag"
	"fmt"
	"github.com/eyebluecn/bamboo/code/core"
	"github.com/eyebluecn/bamboo/code/tool/result"
	"github.com/eyebluecn/bamboo/code/tool/util"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net/http"
	"strings"
	"syscall"
)

const (
	//start web. This is the default mode.
	MODE_WEB     = "web"
	MODE_VERSION = "version"
)

type BambooApplication struct {
	//mode
	mode string

	//EyeblueBamboo host and port  default: http://127.0.0.1:core.DEFAULT_SERVER_PORT
	host string
	//username
	username string
	//password
	password string

	//source file/directory different mode has different usage.
	src string
	//destination directory path(relative to root) in EyeblueBamboo
	dest string
	//true: overwrite, false:skip
	overwrite bool
	filename  string
}

//Start the application.
func (this *BambooApplication) Start() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR:%v\r\n", err)
		}
	}()

	modePtr := flag.String("mode", this.mode, "cli mode web/mirror/crawl")
	hostPtr := flag.String("host", this.username, "bamboo host")
	usernamePtr := flag.String("username", this.username, "username")
	passwordPtr := flag.String("password", this.password, "password")
	srcPtr := flag.String("src", this.src, "src absolute path")
	destPtr := flag.String("dest", this.dest, "destination path in bamboo.")
	overwritePtr := flag.Bool("overwrite", this.overwrite, "whether same file overwrite")
	filenamePtr := flag.String("filename", this.filename, "filename when crawl")

	//flag.Parse() must invoke before use.
	flag.Parse()

	this.mode = *modePtr
	this.host = *hostPtr
	this.username = *usernamePtr
	this.password = *passwordPtr
	this.src = *srcPtr
	this.dest = *destPtr
	this.overwrite = *overwritePtr
	this.filename = *filenamePtr

	//default start as web.
	if this.mode == "" || strings.ToLower(this.mode) == MODE_WEB {

		this.HandleWeb()

	} else if strings.ToLower(this.mode) == MODE_VERSION {

		this.HandleVersion()

	} else {

		//default host.
		if this.host == "" {
			this.host = fmt.Sprintf("http://127.0.0.1:%d", core.DEFAULT_SERVER_PORT)
		}

		if this.username == "" {
			panic(result.BadRequest("in mode %s, username is required", this.mode))
		}

		if this.password == "" {

			if util.EnvDevelopment() {
				panic(result.BadRequest("If run in IDE, use -password yourPassword to input password"))
			} else {

				fmt.Print("Enter Password:")
				bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					panic(err)
				}

				this.password = string(bytePassword)
				fmt.Println()
			}
		}

		panic(result.BadRequest("cannot handle mode %s \r\n", this.mode))
	}

}

func (this *BambooApplication) HandleWeb() {

	//Step 1. Logger
	bambooLogger := &BambooLogger{}
	core.LOGGER = bambooLogger
	bambooLogger.Init()
	defer bambooLogger.Destroy()

	//Step 2. Configuration
	bambooConfig := &BambooConfig{}
	core.CONFIG = bambooConfig
	bambooConfig.Init()

	//Step 3. Global Context
	bambooContext := &BambooContext{}
	core.CONTEXT = bambooContext
	bambooContext.Init()
	defer bambooContext.Destroy()

	//Step 4. Start http
	http.Handle("/", core.CONTEXT)
	core.LOGGER.Info("App started at http://localhost:%v", core.CONFIG.ServerPort())

	dotPort := fmt.Sprintf(":%v", core.CONFIG.ServerPort())
	err1 := http.ListenAndServe(dotPort, nil)
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err1)
	}

}

//fetch the application version
func (this *BambooApplication) HandleVersion() {

	fmt.Printf("EyeblueBamboo %s\r\n", core.VERSION)

}

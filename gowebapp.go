package main

import (
	"log"
	"os"
	"runtime"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/ya42/go_webapp/route"
	"github.com/ya42/go_webapp/common/email"
	"github.com/ya42/go_webapp/common/recaptcha"
	"github.com/ya42/go_webapp/common/server"
	"github.com/ya42/go_webapp/common/session"
	"github.com/ya42/go_webapp/controller"
	"github.com/ya42/go_webapp/plugin"
	"github.com/ya42/go_webapp/service"
)

type Configuration struct {
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  controller.Template   `json:"Template"`
	View      controller.View       `json:"View"`
}

var config = Configuration{}

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	ReadConfig("config"+string(os.PathSeparator)+"config.json")
	session.Configure(config.Session)
	controller.Configure(config.View)
	controller.LoadTemplates(config.Template.Root, config.Template.Children)
	controller.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

func ReadConfig(configFile string) *Configuration{
	var err error
	var absPath string
	var input = io.ReadCloser(os.Stdin)
	if absPath, err = filepath.Abs(configFile); err != nil {
		log.Fatalln(err)
	}
	if input, err = os.Open(absPath); err != nil {
		log.Fatalln(err)
	}
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		log.Fatalln(err)
	}
	service.ParseConfig(jsonBytes, &config)
  return &config
}

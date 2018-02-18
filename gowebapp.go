package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/ya42/go_webapp/route"
	"github.com/ya42/go_webapp/common/email"
	"github.com/ya42/go_webapp/common/jsonconfig"
	"github.com/ya42/go_webapp/common/recaptcha"
	"github.com/ya42/go_webapp/common/server"
	"github.com/ya42/go_webapp/common/session"
	"github.com/ya42/go_webapp/controller"
	"github.com/ya42/go_webapp/plugin"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

	// Setup the views
	controller.Configure(config.View)
	controller.LoadTemplates(config.Template.Root, config.Template.Children)
	controller.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  controller.Template   `json:"Template"`
	View      controller.View       `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

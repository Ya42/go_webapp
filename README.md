# GoWebApp

[![Go Report Card](https://goreportcard.com/badge/github.com/ya42/go_webapp)](https://goreportcard.com/report/github.com/ya42/go_webapp)
[![GoDoc](https://godoc.org/github.com/ya42/go_webapp?status.svg)](https://godoc.org/github.com/ya42/go_webapp) 

Basic MVC Web Application in Go

#### This is a raw MVC framework for GoLang web applications.

This project is inspired and loosely built on Joseph Spurrier's work (https://github.com/josephspurrier/gowebapp).There is also a newer version, Blue Jay(https://github.com/blue-jay/blueprint), by the same author.

This application is built and only tested in Go1.9.1, Ubuntu16.04, bolt1.2.1
No other databases are supported for this version.

To install, run the following command:
~~~
go get github.com/ya42/go_webapp
~~~

## Quick Configuration

All configurations are stored in Config/config.json.
Use port 8080 for http and port 443 for https

## Overview

The web app has a storage layer, a business logic layer, and a presentation layer. The presentation layer is built 
with MVC principles. 

## Structure

The project is organized into the following folders:

~~~
config		- application settings and database schema
static		- location of statically served files like CSS and JS
template	- HTML templates
controller	- page logic organized by HTTP methods (GET, POST)
model		- database object models
route		- route information and middleware
service         - business logic that interfaces datalayer

common		- packages for templates, Boltdatabase, cryptography, sessions
~~~

There are a few external packages I leveraged to create this demo framework:

~~~
github.com/gorilla/context				- registry for global request variables
github.com/gorilla/sessions				- cookie and filesystem sessions
github.com/haisum/recaptcha				- Google reCAPTCHA support
github.com/josephspurrier/csrfbanana 	- CSRF protection for gorilla sessions
github.com/julienschmidt/httprouter 	- high performance HTTP request router
github.com/justinas/alice				- middleware chaining
golang.org/x/crypto/bcrypt 				- password hashing algorithm
~~~

## Templates


## Controllers

The controller files all share the same package name. This cuts down on the 
number of packages when you are mapping the routes. It also forces you to use
a good naming convention for each of the funcs so you know where each of the 
funcs are located and what type of HTTP request they each are mapped to.

### These are a few things you can do with controllers.

A

It's a good idea to abstract the database layer out so if you need to make 
changes, you don't have to look through business logic to find the queries. All
the queries are stored in the models folder.

This project supports BoltDB, MongoDB, and MySQL. All the queries are stored in
the same files so you can easily change the database without modifying anything
but the config file.

The user.go and note.go files are at the root of the model directory and are a
compliation of all the queries for each database type. There are a few hacks in
the models to get the structs to work with all the supported databases.

Connect to the database (only once needed in your application):


## Configuration

To make the web app a little more flexible, you can make changes to different 
components in one place through the config.json file. If you want to add any 
of your own settings, you can add them to config.json and update the structs
in gowebapp.go and the individual files so you can reference them in your code. 
This is config.json:

~~~ json
{
	"Database": {
		"Type": "Bolt",
		"Bolt": {		
 			"Path": "gowebapp.db"
  		},
		"MongoDB": {
			"URL": "127.0.0.1",
			"Database": "gowebapp"
		},
		"MySQL": {
			"Username": "root",
			"Password": "",
			"Name": "gowebapp",
			"Hostname": "127.0.0.1",
			"Port": 3306,
			"Parameter": "?parseTime=true"
		}
	},
	"Email": {
		"Username": "",
		"Password": "",
		"Hostname": "",
		"Port": 25,
		"From": ""
	},
	"Recaptcha": {
		"Enabled": false,
		"Secret": "",
		"SiteKey": ""
	},
	"Server": {
		"Hostname": "",
		"UseHTTP": true,
		"UseHTTPS": false,
		"HTTPPort": 80,
		"HTTPSPort": 443,
		"CertFile": "tls/server.crt",
		"KeyFile": "tls/server.key"
	},
	"Session": {
		"SecretKey": "@r4B?EThaSEh_drudR7P_hub=s#s2Pah",
		"Name": "gosess",
		"Options": {
			"Path": "/",
			"Domain": "",
			"MaxAge": 28800,
			"Secure": false,
			"HttpOnly": true
		}
	},
	"Template": {
		"Root": "base",
		"Children": [
			"partial/menu",
			"partial/footer"
		]
	},
	"View": {
		"BaseURI": "/",
		"Extension": "tmpl",
		"Folder": "template",
		"Name": "blank",
		"Caching": true
	}
}

~~~


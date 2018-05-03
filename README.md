# GoWebApp

[![Go Report Card](https://goreportcard.com/badge/github.com/ya42/go_webapp)](https://goreportcard.com/report/github.com/ya42/go_webapp)
[![GoDoc](https://godoc.org/github.com/ya42/go_webapp?status.svg)](https://godoc.org/github.com/ya42/go_webapp) 

Basic RESTful + MVC Web Application in Go

#### This is a raw RESTful + MVC framework for GoLang web applications.

This project is inspired and loosely built on top of Joseph Spurrier's work (https://github.com/josephspurrier/gowebapp).There is also a newer version, Blue Jay(https://github.com/blue-jay/blueprint), by the same author.

This application is built and only tested in Go 1.9.1, Ubuntu 16.04, bolt database 1.2.1
No other databases are supported for this version.

To install, run the following command:
~~~
go get github.com/ya42/go_webapp
~~~

## Quick Configuration

All configurations are stored in Config/config.json.
Default ports are 8080 for http and port 443 for https

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

There are a few external packages I leveraged to create this framework. Among those I found httprouter and alice particularly useful.

~~~
github.com/gorilla/context				- registry for global request variables
github.com/gorilla/sessions				- cookie and filesystem sessions
github.com/haisum/recaptcha				- Google reCAPTCHA support
github.com/josephspurrier/csrfbanana 	                - CSRF protection for gorilla sessions
github.com/julienschmidt/httprouter 	                - high performance HTTP request router
github.com/justinas/alice				- middleware chaining
golang.org/x/crypto/bcrypt 				- password hashing algorithm
~~~

## Templates
HTML templates. I have not experimented with ajax so far.

## Controllers
Be careful and consistent with naming. Don't forget to synchronize the http route table in route.go

## Database
//todo

## Configuration
config file



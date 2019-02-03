package templates

var (
	applicationTemplate = `package controllers

import (
	"net/http"

	"github.com/dolab/gogo"
	"github.com/dolab/gogo/pkgs/hooks"

	"{{.Namespace}}/{{.Application}}/app/middlewares"
	"{{.Namespace}}/{{.Application}}/app/models"
)

// An Application defines a service meeting gogo.Servicer interface, and wraps
// gogo.Grouper for custom resources.
type Application struct {
	v1 gogo.Grouper
}

// New creates a gogo.Servicer for intialization.
func New() gogo.Servicer {
	return &Application{}
}

// Init implements gogo.Servicer
func (app *Application) Init(config gogo.Configer, group gogo.Grouper) {
	err := NewAppConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// setup models
	err = models.Setup(Config.Model)
	if err != nil {
		panic(err.Error())
	}

	app.v1 = group.NewGroup("/v1")
}

// Filters implements gogo.Servicer
func (app *Application) Filters() {
	// apply your filters for group

	// panic recovery
	app.v1.Use(middlewares.Recovery())
}

// Resources implements gogo.Servicer
func (app *Application) Resources() {
	// register your resources
	// app.v1.GET("/", handler)

	app.v1.GET("/@greeting/hello", GettingStart.Hello)
}

// RequestReceived allows custom request received hooks of server
func (app *Application) RequestReceived() []hooks.NamedHook {
	return []hooks.NamedHook{
		{
			Name: "request_received@debugger",
			Apply: func(w http.ResponseWriter, r *http.Request) bool {
				if Config.Debug {
					log := gogo.NewContextLogger(r)
					log.Debug("RequestReceivedHook")
				}

				return true
			},
		},
	}
}

// RequestRouted allows custom request routed hooks of server
func (app *Application) RequestRouted() []hooks.NamedHook {
	return []hooks.NamedHook{
		{
			Name: "request_routed@debugger",
			Apply: func(w http.ResponseWriter, r *http.Request) bool {
				if Config.Debug {
					log := gogo.NewContextLogger(r)
					log.Debug("RequestRoutedHook")
				}

				return true
			},
		},
	}
}

// ResponseReady allows custom response ready hooks of server
func (app *Application) ResponseReady() []hooks.NamedHook {
	return []hooks.NamedHook{
		{
			Name: "response_ready@debugger",
			Apply: func(w http.ResponseWriter, r *http.Request) bool {
				if Config.Debug {
					log := gogo.NewContextLogger(r)
					log.Debug("ResponseReadyHook")
				}

				return true
			},
		},
	}
}

// ResponseAlways allows custom response always hooks of server
func (app *Application) ResponseAlways() []hooks.NamedHook {
	return []hooks.NamedHook{
		{
			Name: "response_always@debugger",
			Apply: func(w http.ResponseWriter, r *http.Request) bool {
				if Config.Debug {
					log := gogo.NewContextLogger(r)
					log.Debug("ResponseAlwaysHook")
				}

				return true
			},
		},
	}
}
`
)

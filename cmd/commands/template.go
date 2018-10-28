package commands

var (
	envTemplate = `#!/usr/bin/env bash

export GOGOROOT=$(pwd)

# adjust GOPATH
case ":$GOPATH:" in
    *":$GOGOROOT:"*) :;;
    *) GOPATH=$GOPATH:$GOGOROOT;;
esac
export GOPATH


# adjust PATH
readopts="ra"
if [ -n "$ZSH_VERSION" ]; then
    readopts="rA";
fi
while IFS=':' read -$readopts ARR; do
    for i in "${ARR[@]}"; do
        case ":$PATH:" in
            *":$i/bin:"*) :;;
            *) PATH=$i/bin:$PATH
        esac
    done
done <<< "$GOPATH"
export PATH


# mock development && test envs
if [ ! -d "$GOGOROOT/src/{{.Namespace}}/{{.Application}}" ];
then
    mkdir -p "$GOGOROOT/src/{{.Namespace}}"
    ln -s "$GOGOROOT/gogo/" "$GOGOROOT/src/{{.Namespace}}/{{.Application}}"
fi
`

	makefileTemplate = `all: gobuild gotest

godev:
	cd gogo && go run main.go

gobuild: goclean goinstall

gorebuild: goclean goreinstall

goclean:
	go clean ./...

goinstall:
	go get -v github.com/dolab/httpmitm
	go get -v github.com/dolab/httptesting
	go get -v github.com/golib/assert
	go get -v {{.Namespace}}/{{.Application}}

goreinstall:
	go get -v -u github.com/dolab/httpmitm
	go get -v -u github.com/dolab/httptesting
	go get -v -u github.com/golib/assert
	go get -v -u {{.Namespace}}/{{.Application}}

gotest:
	go test {{.Namespace}}/{{.Application}}/app/controllers
	go test {{.Namespace}}/{{.Application}}/app/middlewares
	go test {{.Namespace}}/{{.Application}}/app/models

gopackage:
	mkdir -p bin && go build -a -o bin/{{.Application}} src/{{.Namespace}}/{{.Application}}/main.go

travis: gobuild gotest
`

	gitIgnoreTemplate = `# Compiled Object files, Static and Dynamic libs (Shared Objects)
*.o
*.a
*.so
*.out

# Folders
_obj
_test
bin
pkg
src

# Architecture specific extensions/prefixes
*.[568vq]
[568vq].out

*.cgo1.go
*.cgo2.c
_cgo_defun.c
_cgo_gotypes.go
_cgo_export.*
_testmain.go

*.exe
*.test
*.prof

# development & test config files
*.development.json
*.test.json
`

	mainTemplate = `package main

import (
	"flag"
	"os"
	"path"

	"github.com/dolab/gogo"

	"{{.Namespace}}/{{.Application}}/app/controllers"
)

var (
	runMode string // app run mode, available values are [development|test|production], default to development
	srcPath string // app config path, e.g. /home/deploy/websites/helloapp
)

func main() {
	flag.StringVar(&runMode, "runMode", "development", "{{.Application}} -runMode=[development|test|production]")
	flag.StringVar(&srcPath, "srcPath", "", "{{.Application}} -srcPath=/path/to/[config/application.json]")
	flag.Parse()

	// verify run mode
	if mode := gogo.RunMode(runMode); !mode.IsValid() {
		flag.PrintDefaults()
		return
	}

	// adjust src path
	if srcPath == "" {
		var err error

		srcPath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		srcPath = path.Clean(srcPath)
	}

	controllers.New(runMode, srcPath).Run()
}
`

	applicationTemplates = map[string]string{
		"application": `package controllers

import (
	"github.com/dolab/gogo"

	"{{.Namespace}}/{{.Application}}/app/middlewares"
	"{{.Namespace}}/{{.Application}}/app/models"
)

// Application extends gogo.AppServer by customization
type Application struct {
	*gogo.AppServer
}

func New(runMode, srcPath string) *Application {
	appServer := gogo.New(runMode, srcPath)

	// init Config
	err := NewAppConfig(appServer.Config())
	if err != nil {
		panic(err.Error())
	}

	// init models
	err = models.Setup(Config.Model)
	if err != nil {
		panic(err.Error())
	}

	return &Application{appServer}
}

// Middlerwares implements gogo.Middlewarer
// NOTE: DO NOT change the method name, its required by gogo!
func (app *Application) Middlewares() {
	// apply your middlewares

	// panic recovery
	app.Use(middlewares.Recovery())
}

// Resources implements gogo.Resourcer
// NOTE: DO NOT change the method name, its required by gogo!
func (app *Application) Resources() {
	// register your resources
	// app.GET("/", handler)

	app.GET("/@getting_start/hello", GettingStart.Hello)
}

// Run runs application after registering middelwares and resources
func (app *Application) Run() {
	// register middlewares
	app.Middlewares()

	// register resources
	app.Resources()

	// run server
	app.AppServer.Run()
}
`,
		"application_testing": `package controllers

import (
	"os"
	"path"
	"testing"

	"github.com/dolab/httptesting"
)

var (
	gogotest *httptesting.Client
)

func TestMain(m *testing.M) {
	var (
		runMode = "test"
		srcPath = path.Clean("../../")
	)

	app := New(runMode, srcPath)
	app.Resources()

	gogotest = httptesting.NewServer(app, false)

	code := m.Run()

	gogotest.Close()

	os.Exit(code)
}
`,
		"application_config": `package controllers

import (
	"github.com/dolab/gogo"

	"{{.Namespace}}/{{.Application}}/app/models"
)

var (
	Config *AppConfig
)

// AppConfig defines specs for application config
type AppConfig struct {
	Model        *models.Config      ` + "`" + `json:"model"` + "`" + `
	Domain       string              ` + "`" + `json:"domain"` + "`" + `
	GettingStart *GettingStartConfig ` + "`" + `json:"getting_start"` + "`" + `
}

// NewAppConfig apply application config from gogo.Configer
func NewAppConfig(config gogo.Configer) error {
	return config.UnmarshalJSON(&Config)
}

// Sample application config for illustration
type GettingStartConfig struct {
	Greeting string ` + "`" + `json:"greeting"` + "`" + `
}
`,
		"application_config_test": `package controllers

import (
	"testing"

	"github.com/golib/assert"
)

func Test_AppConfig(t *testing.T) {
	assertion := assert.New(t)

	assertion.NotEmpty(Config.Domain)
	assertion.NotNil(Config.GettingStart)
}
`,
		"getting_start": `package controllers

import (
	"github.com/dolab/gogo"
)

var (
	GettingStart *_GettingStart
)

type _GettingStart struct{}

// @route GET /@getting_start/hello
func (_ *_GettingStart) Hello(ctx *gogo.Context) {
	name := ctx.Params.Get("name")
	if name == "" {
		name = Config.GettingStart.Greeting
	}

	ctx.Text(name)
}
`,
		"getting_start_test": `package controllers

import (
	"net/url"
	"testing"
)

func Test_GettingStart_Hello(t *testing.T) {
	// it should work without greeting
	request := gogotest.New(t)
	request.Get("/@getting_start/hello")

	request.AssertOK()
	request.AssertContains(Config.GettingStart.Greeting)

	// it should work with custom greeting
	greeting := "Hi, gogo!"

	params := url.Values{}
	params.Add("name", greeting)

	request = gogotest.New(t)
	request.Get("/@getting_start/hello", params)

	request.AssertOK()
	request.AssertContains(greeting)
}
`,
		"middleware_testing": `package middlewares

import (
	"os"
	"path"
	"testing"

	"github.com/dolab/gogo"
	"github.com/dolab/httptesting"
)

var (
	gogoapp  *gogo.AppServer
	gogotest *httptesting.Client
)

func TestMain(m *testing.M) {
	var (
		runMode = "test"
		srcPath = path.Clean("../../")
	)

	gogoapp = gogo.New(runMode, srcPath)
	gogotest = httptesting.NewServer(gogoapp, false)

	code := m.Run()

	gogotest.Close()

	os.Exit(code)
}
`,
		"middleware_recovery": `package middlewares

import (
	"runtime"
	"strings"

	"github.com/dolab/gogo"
)

func Recovery() gogo.Middleware {
	return func(ctx *gogo.Context) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				// where does panic occur? try max 20 depths
				pcs := make([]uintptr, 20)
				max := runtime.Callers(2, pcs)

				if max == 0 {
					ctx.Logger.Warn("No pcs available.")
				} else {
					frames := runtime.CallersFrames(pcs[:max])
					for {
						frame, more := frames.Next()

						// To keep this example's output stable
						// even if there are changes in the testing package,
						// stop unwinding when we leave package runtime.
						if strings.Contains(frame.Function, "runtime.") {
							if more {
								continue
							} else {
								break
							}
						}

						tmp := strings.SplitN(frame.File, "/src/", 2)
						if len(tmp) == 2 {
							ctx.Logger.Errorf("(src/%s:%d: %v)", tmp[1], frame.Line, panicErr)
						} else {
							ctx.Logger.Errorf("(%s:%d: %v)", frame.File, frame.Line, panicErr)
						}

						break
					}
				}

				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}
`,
		"middleware_recovery_test": `package middlewares

import (
	"testing"

	"github.com/dolab/gogo"
)

func Test_Recovery(t *testing.T) {
	// register temp resource for testing
	app := gogoapp.Group("", Recovery())

	app.GET("/middlewares/recovery", func(ctx *gogo.Context) {
		panic("Recover testing")
	})

	// it should work
	request := gogotest.New(t)
	request.Get("/middlewares/recovery")

	request.AssertOK()
}
`,
		"model": `package models

import (
	"database/sql"
)

var (
	model *sql.Conn
)

// Config defines config for model driver
type Config struct {
	Host     string ` + "`" + `json:"host"` + "`" + `
	Username string ` + "`" + `json:"username"` + "`" + `
	Password string ` + "`" + `json:"password"` + "`" + `
}

// Setup inject model with driver conn
func Setup(config *Config) error {
	// TODO: create *sql.Conn with config

	model = &sql.Conn{}

	return nil
}
`,
		"model_test": `package models

import (
	"testing"

	"github.com/golib/assert"
)

func Test_Setup(t *testing.T) {
	assertion := assert.New(t)

	assertion.Nil(model)
	Setup(&Config{})
	assertion.NotNil(model)
}
`,
		"application_config_json": `{
	"name": "{{.Application}}",
	"mode": "test",
	"sections": {
		"development": {
			"server": {
				"addr": "localhost",
				"port": 9090,
				"ssl": false,
				"request_timeout": 30,
				"response_timeout": 30,
				"request_id": "X-Request-Id"
			},
			"logger": {
				"output": "stdout",
				"level": "debug",
				"filter_params": ["password", "password_confirmation"]
			},
			"domain": "https://example.com",
			"getting_start": {
				"greeting": "Hello, gogo!"
			}
		},

		"test": {
			"server": {
				"addr": "localhost",
				"port": 9090,
				"ssl": false,
				"request_timeout": 30,
				"response_timeout": 30,
				"request_id": "X-Request-Id"
			},
			"logger": {
				"output": "stdout",
				"level": "info",
				"filter_params": ["password", "password_confirmation"]
			},
			"domain": "https://example.com",
			"getting_start": {
				"greeting": "Hello, gogo!"
			}
		},

		"production": {
			"server": {
				"addr": "localhost",
				"port": 9090,
				"ssl": true,
				"ssl_cert": "/path/to/ssl/cert",
				"ssl_key": "/path/to/ssl/key",
				"request_timeout": 30,
				"response_timeout": 30,
				"request_id": "X-Request-Id"
			},
			"logger": {
				"output": "stdout",
				"level": "warn",
				"filter_params": ["password", "password_confirmation"]
			}
		}
	}
}
`}

	componentTemplates = map[string]string{
		"controller": `package controllers

import (
	"net/http"

	"github.com/dolab/gogo"
)

var (
	{{.Name}} *_{{.Name}}
)

type _{{.Name}} struct{}

// // custom resource id name, default to {{.Name|lowercase}}
// func (_ *_{{.Name}}) ID() string {
// 	return "id"
// }

// @route GET /{{.Name|lowercase}}
func (_ *_{{.Name}}) Index(ctx *gogo.Context) {
	ctx.SetStatus(http.StatusNotImplemented)
	ctx.Return()
}

// @route POST /{{.Name|lowercase}}
func (_ *_{{.Name}}) Create(ctx *gogo.Context) {
	ctx.SetStatus(http.StatusNotImplemented)
	ctx.Return()
}

// @route GET /{{.Name|lowercase}}/:{{.Name|lowercase}}
func (_ *_{{.Name}}) Show(ctx *gogo.Context) {
	// retrieve resource name of path params
	id := ctx.Params.Get("{{.Name|lowercase}}")

	ctx.SetStatus(http.StatusNotImplemented)
	ctx.Json(map[string]interface{}{
		"id": id,
		"tags": []string{
			id,
		},
	})
}

// @route PUT /{{.Name|lowercase}}/:{{.Name|lowercase}}
func (_ *_{{.Name}}) Update(ctx *gogo.Context) {
	// retrieve resource name of path params
	id := ctx.Params.Get("{{.Name|lowercase}}")

	ctx.SetStatus(http.StatusNotImplemented)
	ctx.Return(id)
}

// @route DELETE /{{.Name|lowercase}}/:{{.Name|lowercase}}
func (_ *_{{.Name}}) Destroy(ctx *gogo.Context) {
	// retrieve resource name of path params
	id := ctx.Params.Get("{{.Name|lowercase}}")

	ctx.SetStatus(http.StatusNotImplemented)
	ctx.Return(id)
}
`,
		"controller_test": `package controllers

import (
	"net/http"
	"testing"
)

func Test_{{.Name}}_Index(t *testing.T) {
	request := gogotest.New(t)
	request.Get("/{{.Name|lowercase}}")

	request.AssertStatus(http.StatusNotImplemented)
	request.AssertEmpty()
}

func Test_{{.Name}}_Create(t *testing.T) {
	request := gogotest.New(t)
	request.PostJSON("/{{.Name|lowercase}}", nil)

	request.AssertStatus(http.StatusNotImplemented)
	request.AssertEmpty()
}

func Test_{{.Name}}_Show(t *testing.T) {
	id := "{{.Name|lowercase}}"

	request := gogotest.New(t)
	request.Get("/{{.Name|lowercase}}/" + id)

	request.AssertStatus(http.StatusNotImplemented)
	request.AssertContainsJSON("id", id)
	request.AssertContainsJSON("tags.0", id)
}

func Test_{{.Name}}_Update(t *testing.T) {
	id := "{{.Name|lowercase}}"

	request := gogotest.New(t)
	request.PutJSON("/{{.Name|lowercase}}/"+id, nil)

	request.AssertStatus(http.StatusNotImplemented)
	request.AssertContains(id)
}

func Test_{{.Name}}_Destroy(t *testing.T) {
	id := "{{.Name|lowercase}}"

	request := gogotest.New(t)
	request.DeleteJSON("/{{.Name|lowercase}}/"+id, nil)

	request.AssertStatus(http.StatusNotImplemented)
	request.AssertContains(id)
}
`,
		"middleware": `package middlewares

import (
	"github.com/dolab/gogo"
)

func {{.Name}}() gogo.Middleware {
	return func(ctx *gogo.Context) {
		// TODO: implements custom logic
		ctx.AddHeader("x-gogo-middleware", "Hello, middleware!")

		ctx.Next()
	}
}
`,
		"middleware_test": `package middlewares

import (
	"testing"

	"github.com/dolab/gogo"
)

func Test_{{.Name}}(t *testing.T) {
	// register temp resource for testing
	app := gogoapp.Group("", {{.Name}}())

	app.GET("/middlewares/{{.Name|lowercase}}", func(ctx *gogo.Context) {
		ctx.Return()
	})

	request := gogotest.New(t)
	request.Get("/middlewares/{{.Name|lowercase}}", nil)
	request.AssertOK()
	request.AssertHeader("x-gogo-middleware", "Hello, middleware!")
}
`,
		"model": `package models

import (
	"errors"
)

var (
	{{.Name}} *_{{.Name}}
)

type {{.Name}}Model struct {
	// TODO: fill with table fields
}

type _{{.Name}} struct{}

func (_ *_{{.Name}}) Find(id string) (m *{{.Name}}Model, err error) {
	err = errors.New("Not Found")

	return
}
`,
		"model_test": `package models

import (
	"testing"

	"github.com/golib/assert"
)

func Test_{{.Name}}Model(t *testing.T) {
	assertion := assert.New(t)

	m := &{{.Name}}Model{}
	assertion.NotNil(m)
}

func Test_{{.Name}}_Find(t *testing.T) {
	assertion := assert.New(t)

	id := "???"

	m, err := {{.Name}}.Find(id)
	assertion.EqualError(err, "Not Found")
	assertion.Nil(m)
}
`}
)

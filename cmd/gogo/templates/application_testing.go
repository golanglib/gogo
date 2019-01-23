package templates

var (
	applicationTestingTemplate = `package controllers

import (
	"os"
	"path"
	"testing"

	"github.com/dolab/gogo"
	"github.com/dolab/httptesting"
)

var (
	gogotesting *httptesting.Client
)

func TestMain(m *testing.M) {
	var (
		runMode = "test"
		srcPath = path.Clean("../../")
	)

	app := gogo.New(runMode, srcPath)
	app.NewResources(New())

	gogotesting = httptesting.NewServer(app, false)

	code := m.Run()

	gogotesting.Close()

	os.Exit(code)
}
`
)

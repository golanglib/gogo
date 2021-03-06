package templates

import (
	"regexp"
	"strings"
	"text/template"
)

var (
	box *template.Template
)

func Box() *template.Template {
	return box
}

func init() {
	// register templates
	box = template.New("gogo").Funcs(template.FuncMap{
		"lowercase": strings.ToLower,
	})

	// commons
	template.Must(box.New("gitignore").Parse(gitIgnoreTemplate))
	template.Must(box.New("env.sh").Parse(envTemplate))
	template.Must(box.New("go.mod").Parse(modTemplate))
	template.Must(box.New("readme").Parse(
		format(readmeTemplate),
	))
	template.Must(box.New("makefile").Parse(
		format(makefileTemplate),
	))
	template.Must(box.New("main.yml").Parse(
		format(mainConfigTemplate),
	))
	template.Must(box.New("main.go").Parse(
		format(mainTemplate),
	))
	template.Must(box.New("errors.go").Parse(
		format(errorsTemplate),
	))
	// controllers
	template.Must(box.New("application.go").Parse(
		format(applicationTemplate),
	))
	template.Must(box.New("application_testing.go").Parse(
		format(applicationTestingTemplate),
	))
	template.Must(box.New("application_config.go").Parse(
		format(applicationConfigTemplate),
	))
	template.Must(box.New("application_config_test.go").Parse(
		format(applicationConfigTestTemplate),
	))
	template.Must(box.New("application_config.yml").Parse(
		formatYAML(applicationConfigYAMLTemplate),
	))
	template.Must(box.New("getting_start.go").Parse(
		format(gettingStartTemplate),
	))
	template.Must(box.New("getting_start_test.go").Parse(
		format(gettingStartTestTemplate),
	))
	// middlewares
	template.Must(box.New("middleware_testing.go").Parse(
		format(middlewareTestingTemplate),
	))
	template.Must(box.New("middleware_recovery.go").Parse(
		format(middlewareRecoveryTemplate),
	))
	template.Must(box.New("middleware_recovery_test.go").Parse(
		format(middlewareRecoveryTestTemplate),
	))
	// models
	template.Must(box.New("model.go").Parse(
		format(modelTemplate),
	))
	template.Must(box.New("model_test.go").Parse(
		format(modelTestTemplate),
	))
	// templates
	template.Must(box.New("template_controller").Parse(
		format(componentControllerTemplate),
	))
	template.Must(box.New("template_controller_test").Parse(
		format(componentControllerTestTemplate),
	))
	template.Must(box.New("template_filter").Parse(
		format(componentFilterTemplate),
	))
	template.Must(box.New("template_filter_test").Parse(
		format(componentFilterTestTemplate),
	))
	template.Must(box.New("template_model").Parse(
		format(componentModelTemplate),
	))
	template.Must(box.New("template_model_test").Parse(
		format(componentModelTestTemplate),
	))
}

var (
	langOpenTag  = regexp.MustCompile(`<(\w+)>`)
	langCloseTag = regexp.MustCompile(`</(\w+)>`)
	langTab      = regexp.MustCompile(`(?m)^(\t{1,})`)
)

func format(tpl string) string {
	tpl = langOpenTag.ReplaceAllStringFunc(tpl, func(tag string) string {
		return "```" + strings.Trim(tag, "<>")
	})
	tpl = langCloseTag.ReplaceAllStringFunc(tpl, func(tag string) string {
		return "```"
	})

	return tpl
}

func formatYAML(tpl string) string {
	tpl = format(tpl)

	tpl = langTab.ReplaceAllStringFunc(tpl, func(tabs string) string {
		return strings.Repeat("    ", len(tabs))
	})

	return tpl
}

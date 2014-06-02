package html

import (
	"html/template"
	"strings"
)

var (
	defaultFuncs = map[string]interface{}{

		// Skips sanitation on the parameter.
		"raw": func(text string) template.HTML {
			return template.HTML(text)
		},

		// Replaces newlines with '<br>'.
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
	}
)

func createFuncMap(ldr *Loader) map[string]interface{} {
	runView := func(view *View, data ...interface{}) (template.HTML, error) {
		html, err := view.HTML(data...)
		if err != nil {
			return "", err
		}
		return template.HTML(html), nil
	}

	runSet := func(set *Set, data ...interface{}) (template.HTML, error) {
		view, err := ldr.NewSet().AddSet(set).View()
		if err != nil {
			return "", err
		}
		return runView(view, data...)
	}

	runTemplate := func(name string, data ...interface{}) (template.HTML, error) {
		view, err := ldr.NewSet().Add(name).View()
		if err != nil {
			return "", err
		}
		return runView(view, data...)
	}

	funcs := make(map[string]interface{}, len(defaultFuncs)+3)
	funcs["runSet"] = runSet
	funcs["runView"] = runView
	funcs["runTemplate"] = runTemplate
	for name, fn := range defaultFuncs {
		funcs[name] = fn
	}
	return funcs
}

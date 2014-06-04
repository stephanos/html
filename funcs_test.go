package html

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/strutil"
)

var _ = Describe("Template Funcs", func() {

	var loader *Loader
	var reverseFn = func(text string) string { return strutil.Reverse(text) }

	BeforeEach(func() {
		loader = fixtureLoader(false)
	})

	It("function raw", func() {
		out, err := loader.NewSet().Add("funcs/raw").ViewMust().HTML()

		Check(err, IsNil)
		Check(out, Equals, "<br>")
	})

	It("function nl2br", func() {
		out, err := loader.NewSet().Add("funcs/nl2br").ViewMust().HTML()

		Check(err, IsNil)
		Check(out, Equals, "<br>")
	})

	It("function runSet", func() {
		set := loader.NewSet().Add("pages/nonsense")
		_, err := loader.NewSet().Add("funcs/run_set").ViewMust().HTML(set)

		Check(err, NotNil).And(Contains, `error calling runSet: template "pages/nonsense" not found`)
	})

	It("function runView", func() {
		view := loader.NewSet().Add("pages/content").ViewMust()
		out, err := loader.NewSet().Add("funcs/run_view").ViewMust().HTML(view)

		Check(err, IsNil)
		Check(out, Equals, "<h1>Content</h1>")
	})

	It("function runView: error", func() {
		view := loader.NewSet()
		_, err := loader.NewSet().Add("funcs/run_view").ViewMust().HTML(view)

		Check(err, NotNil).And(Contains, `wrong type for value; expected *html.View; got *html.Set`)
	})

	It("function runTemplate", func() {
		tmpl := "pages/content"
		out, err := loader.NewSet().Add("funcs/run_template").ViewMust().HTML(tmpl)

		Check(err, IsNil)
		Check(out, Equals, "<h1>Content</h1>")
	})

	It("function runTemplate: error", func() {
		tmpl := "pages/nonsense"
		_, err := loader.NewSet().Add("funcs/run_template").ViewMust().HTML(tmpl)

		Check(err, NotNil).And(Contains, `error calling runTemplate: template "pages/nonsense" not found`)
	})

	It("add single func", func() {
		out, err := loader.NewSet().Add("funcs/dynamic1").AddFunc("customFunc", reverseFn).ViewMust().HTML()

		Check(err, IsNil)
		Check(out, Equals, "cba")
	})

	It("add multiple funcs", func() {
		out, err := loader.NewSet().
			Add("funcs/dynamic2").
			AddFuncs(map[string]interface{}{"customFunc1": reverseFn, "customFunc2": reverseFn}).
			ViewMust().
			HTML()

		Check(err, IsNil)
		Check(trim(out), Equals, "abc")
	})
})

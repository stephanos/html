package html

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Template Funcs", func() {

	var loader *Loader

	BeforeEach(func() {
		loader = fixtureLoader(false)
	})

	It("raw", func() {
		out, err := loader.NewSet().Add("funcs/raw").ViewMust().HTML()

		Check(err, IsNil)
		Check(out, Equals, "<br>")
	})

	It("nl2br", func() {
		out, err := loader.NewSet().Add("funcs/nl2br").ViewMust().HTML()

		Check(err, IsNil)
		Check(out, Equals, "<br>")
	})

	It("runSet", func() {
		set := loader.NewSet().Add("pages/content")
		out, err := loader.NewSet().Add("funcs/run_set").ViewMust().HTML(set)

		Check(err, IsNil)
		Check(out, Equals, "<h1>Content</h1>")
	})

	It("runView", func() {
		view := loader.NewSet().Add("pages/content").ViewMust()
		out, err := loader.NewSet().Add("funcs/run_view").ViewMust().HTML(view)

		Check(err, IsNil)
		Check(out, Equals, "<h1>Content</h1>")
	})

	It("runTemplate", func() {
		tmpl := "pages/content"
		out, err := loader.NewSet().Add("funcs/run_template").ViewMust().HTML(tmpl)

		Check(err, IsNil)
		Check(out, Equals, "<h1>Content</h1>")
	})
})

package html

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("HTML", func() {

	var loader *Loader

	BeforeEach(func() {
		loader = fixtureLoader(false)
	})

	It("render: add template", func() {
		view, err := loader.NewSet().Add("layout", "pages/home").View()
		Check(err, IsNil)
		Check(view, NotNil)

		out, err := view.HTML()
		Check(err, IsNil)
		Check(trim(out), Equals, `<html> <body> <h1>Home</h1> </body> </html>`)
	})

	It("render: add & set template", func() {
		view, err := loader.NewSet().Add("layout").Set("content", "pages/content").View()
		Check(err, IsNil)
		Check(view, NotNil)

		out, err := view.HTML()
		Check(err, IsNil)
		Check(trim(out), Equals, `<html> <body> <h1>Content</h1> </body> </html>`)
	})

	It("render: add set", func() {
		baseSet := loader.NewSet().Add("layout")
		view, err := loader.NewSet().AddSet(baseSet).Add("pages/home").View()

		out, err := view.HTML()
		Check(err, IsNil)
		Check(trim(out), Equals, `<html> <body> <h1>Home</h1> </body> </html>`)
	})

	It("reload template source", func() {
		loader = fixtureLoader(true)
		view := loader.NewSet().Add("dynamic").ViewMust()

		out, err := view.HTML()
		Check(err, IsNil)
		Check(trim(out), IsEmpty)

		setDynamicTemplate("dynamic")

		out, err = view.HTML()
		Check(err, IsNil)
		Check(trim(out), Equals, "dynamic")

		setDynamicTemplate("{{ invalid }}")

		out, err = view.HTML()
		Check(err, NotNil).And(Contains, `function "invalid" not defined`)
	})

	It("error for missing template", func() {
		_, err := loader.NewSet().Add("not-existing").View()
		Check(err, NotNil).And(Contains, `template "not-existing" not found`)
	})

	It("error for invalid syntax", func() {
		_, err := loader.NewSet().Add("invalid_syntax").View()
		Check(err, NotNil).And(Contains, `unexpected unterminated quoted string in command`)
	})

	It("error for missing function", func() {
		_, err := loader.NewSet().Add("invalid_func").View()
		Check(err, NotNil).And(Contains, `function "invalid" not defined`)
	})

	It("error for root template redefinition", func() {
		_, err := loader.NewSet().Add("layout", "layout").View()
		Check(err, NotNil).And(Contains, `redefinition of root template`)
	})

	It("error for incomplete template", func() {
		_, err := loader.NewSet().Add("pages/home").View()
		Check(err, NotNil).And(Contains, "missing root template")

		_, err = loader.NewSet().Add("layout").View()
		Check(err, NotNil).And(Contains, `missing template(s) ["content"]`)
	})

	It("panic for error", func() {
		Check(func() {
			loader.NewSet().Add("not-existing").ViewMust()
		}, Panics)
	})
})

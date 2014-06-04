package html

import . "github.com/101loops/bdd"

var _ = Describe("Loader", func() {

	It("load template files", func() {
		conf := Config{Directories: []string{"fixtures"}}
		loader, err := NewLoader(conf)

		Check(err, IsNil)
		Check(loader, NotNil)
		Check(loader.Sources(), HasLen, 14)

		loader.AddFile("my-layout", "layout.html")
		loader.AddText("my-content", "<h1>My Content</h1>")

		Check(loader.Sources(), HasLen, 16)
		Check(loader.Sources(), Contains, &Source{Name: "my-layout", FilePath: "layout.html"})
		Check(loader.Sources(), Contains, &Source{Name: "my-content", Content: "<h1>My Content</h1>"})
	})

	It("unable to load non-existent directory", func() {
		conf := Config{Directories: []string{"nonsense"}}
		_, err := NewLoader(conf)

		Check(err, NotNil)
	})
})

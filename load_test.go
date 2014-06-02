package html

import . "github.com/101loops/bdd"

var _ = Describe("Loader", func() {

	It("load template files", func() {
		conf := Config{Directories: []string{"fixtures"}}
		loader, err := NewLoader(conf)

		Check(err, IsNil)
		Check(loader, NotNil)
		Check(loader.Sources(), HasLen, 12)
	})

	It("unable to load non-existent directory", func() {
		conf := Config{Directories: []string{"nonsense"}}
		_, err := NewLoader(conf)

		Check(err, NotNil)
	})
})

package html

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/strutil"
	"io/ioutil"
	"testing"
)

func TestSuite(t *testing.T) {
	setDynamicTemplate("")
	RunSpecs(t, "HTML Suite")
	setDynamicTemplate("")
}

func setDynamicTemplate(content string) {
	err := ioutil.WriteFile("fixtures/dynamic.html", []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

func fixtureLoader(reload bool) *Loader {
	conf := Config{Directories: []string{"fixtures"}, AutoReload: reload}
	loader, err := NewLoader(conf)
	if err != nil {
		panic(err)
	}
	return loader
}

func trim(out string) string {
	return strutil.ShrinkWhitespace(out)
}

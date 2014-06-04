package html

import (
	"bytes"
	"html/template"
	"io"
)

// A View can be rendered to HTML text.
type View struct {
	set      *Set
	template *template.Template
}

// Write applies the parsed template to the specified data object,
// writing the output to wr.
func (v *View) Write(wr io.Writer, data ...interface{}) error {
	err := v.set.create(v)
	if err != nil {
		return err
	}

	var input interface{}
	if len(data) != 0 {
		input = data[0]
	}

	return v.template.ExecuteTemplate(wr, rootTemplateName, input)
}

// HTML applies the parsed template to the specified data object,
// returning the result as a text.
func (v *View) HTML(data ...interface{}) (string, error) {
	var doc bytes.Buffer
	err := v.Write(&doc, data...)
	return doc.String(), err
}

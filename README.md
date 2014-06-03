html [![Build Status](https://secure.travis-ci.org/101loops/html.png)](https://travis-ci.org/101loops/html) [![Coverage Status](https://coveralls.io/repos/101loops/html/badge.png)](https://coveralls.io/r/101loops/html) [![GoDoc](https://camo.githubusercontent.com/6bae67c5189d085c05271a127da5a4bbb1e8eb2c/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f736d61727479737472656574732f676f636f6e7665793f7374617475732e706e67)](http://godoc.org/github.com/101loops/html)
======

This Go package to load, compose and render HTML templates.


### Features

- **fluent API:** easily compose template into sets
- **auto-reloading:** reload templates on page refresh
- **redefinition:** define a default and overwrite it later 
- **validation:** ensure completeness at time of creation (not rendering) 
- **helper functions:** e.g. use 'runTemplate' to execute an arbitrary template
- **caching:** templates are only parsed once


### Example

```go
import "github.com/101loops/html"

// specify template source directories, enable auto-reload
conf := html.Config{Directories: []string{"views"}, AutoReload: true}

// scan for available templates
loader, _ := html.NewLoader(conf)

// create two sets: a re-usable and a specific one
baseSet := loader.NewSet().Add("layout", "partials/header", "partials/footer")
helloSet := loader.NewSet().AddSet(baseSet).Add("pages/hello")

// create executable view, making sure all template placeholders are defined
view := helloSet.ViewMust()

// execute the template and print the result to a Writer
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    view.Write(w, "World") 
})
```


### ToDos
 
- render any error to HTML (+ display snippet of template source)
- allow custom template parser


### Installation
`go get github.com/101loops/html`

### Documentation
[godoc.org](http://godoc.org/github.com/101loops/html)

### License
MIT (see LICENSE).

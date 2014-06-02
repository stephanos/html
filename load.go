package html

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template/parse"
)

// Loader collects the available template sources.
// It creates new sets and triggers the template parsing.
type Loader struct {
	conf      Config
	mutex     sync.Mutex
	sources   map[string]*Source
	treeCache map[string]map[string]*parse.Tree
}

// Config controls the behaviour of the template loading and rendering.
type Config struct {

	// Directories define the file paths to search for templates,
	// ordered by descending priority.
	Directories []string

	// AutoReload is whether the templates should be reloaded on rendering.
	// This is useful in development, but should be disabled in production.
	AutoReload bool

	// DelimLeft is the delimiter that marks the start of a template action.
	DelimLeft string

	// DelimRight is the delimiter that marks the stop of a template action.
	DelimRight string
}

// Source is a template data source. It contains an optional path to
// the source's file source.
type Source struct {
	Content  string
	FilePath string
	Name     string
}

// NewLoader creates a new loader. It scans the provided source directories
// and collects all available template sources.
func NewLoader(conf Config) (l *Loader, err error) {
	l = &Loader{conf: conf, treeCache: make(map[string]map[string]*parse.Tree)}
	err = l.scan()
	return l, err
}

func (l *Loader) scan() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	directories := l.conf.Directories
	l.sources = make(map[string]*Source)
	for i := len(l.conf.Directories); i > 0; i-- {
		dir := directories[i-1]

		if _, err := os.Stat(dir); err != nil {
			return err
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, fileErr error) error {
			if fileErr != nil {
				return fileErr
			}

			fileName := info.Name()
			ignoreFile := strings.HasPrefix(fileName, "_") || !strings.HasSuffix(fileName, ".html")
			if info.IsDir() || ignoreFile {
				return nil
			}

			viewName, fileErr := filepath.Rel(dir, path)
			if fileErr != nil {
				return fileErr
			}
			if os.PathSeparator == '\\' {
				viewName = strings.Replace(viewName, `\`, `/`, -1) // replaces path separator on windows
			}
			viewName = strings.Replace(viewName, ".html", "", 1)

			l.sources[viewName] = &Source{Name: viewName, FilePath: path}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// NewSet returns a new initialized set.
func (l *Loader) NewSet() *Set {
	return newSet(l)
}

// AddFile adds a file-based template source.
// It overwrites any existing source with the same name.
func (l *Loader) AddFile(name, filePath string) *Loader {
	return l.addSource(&Source{Name: name, FilePath: filePath})
}

// AddText adds a text-based template source.
// It overwrites any existing source with the same name.
func (l *Loader) AddText(name string, content string) *Loader {
	return l.addSource(&Source{Name: name, Content: content})
}

func (l *Loader) addSource(source *Source) *Loader {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.sources[source.Name] = source
	return l
}

// Sources returns the sources found by scanning the provided directories.
func (l *Loader) Sources() []*Source {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	ret := make([]*Source, len(l.sources))
	var i int
	for _, src := range l.sources {
		ret[i] = src
		i++
	}
	return ret
}

func (l *Loader) buildTemplate(s *Set) (*template.Template, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	err := l.loadSources(s)
	if err != nil {
		return nil, err
	}

	tmpl, err := l.parseSources(s)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (l *Loader) loadSources(s *Set) error {
	for _, source := range s.sources {
		var loadSourceFromPath string
		if source.FilePath != "" {
			if src, ok := l.sources[source.FilePath]; ok {
				loadSourceFromPath = src.FilePath
			} else {
				return fmt.Errorf("template %q not found", source.FilePath)
			}
		}

		content, err := ioutil.ReadFile(loadSourceFromPath)
		if err != nil {
			return err
		}
		source.Content = string(content)
	}
	return nil
}

func (l *Loader) parseSources(s *Set) (*template.Template, error) {
	var useCache = !l.conf.AutoReload
	var parsedTrees []*parse.Tree
	parsedNames := make(map[string]bool)

	for _, source := range s.sources {
		var trees map[string]*parse.Tree

		cacheKey := source.FilePath
		if useCache && cacheKey != "" {
			if cachedTrees, ok := l.treeCache[cacheKey]; ok {
				trees = cachedTrees
			}
		}

		var err error
		trees, err = parse.Parse(source.Name, source.Content, l.conf.DelimLeft, l.conf.DelimRight, s.funcs)
		if err != nil {
			return nil, err
		}

		for name, tree := range trees {
			if name == rootTemplateName {
				isRootTmpl := len(trees) == 1 && trees[rootTemplateName] != nil
				if !isRootTmpl {
					continue
				}

				alreadyParsed := parsedNames[name]
				if alreadyParsed {
					return nil, fmt.Errorf("html/template: redefinition of root template")
				}
			}

			parsedTrees = append(parsedTrees, tree)
			parsedNames[name] = true
		}

		if useCache && cacheKey != "" {
			l.treeCache[cacheKey] = trees
		}
	}

	return createTemplate(parsedTrees, s.funcs)
}

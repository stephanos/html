package html

// Set is a collection of template sources.
// It allows to create a View from its sources.
type Set struct {
	loader  *Loader
	sources []*Source
	funcs   map[string]interface{}
}

func newSet(l *Loader) *Set {
	set := &Set{
		loader: l,
		funcs:  createFuncMap(l),
	}
	return set
}

// Add appends a template source, identified by its file path, to the set.
// A copy of the set is returned.
func (s *Set) Add(files ...string) *Set {
	ret := s.clone()
	for _, file := range files {
		ret.sources = append(ret.sources, &Source{FilePath: file})
	}
	return ret
}

// AddSet appends the template sources of another set to this one.
// A copy of the set is returned.
func (s *Set) AddSet(set *Set) *Set {
	return s.AddSets(set)
}

// AddSets appends the template sources of other sets to this one.
// It returns a copy of the Set.
func (s *Set) AddSets(sets ...*Set) *Set {
	ret := s.clone()
	for _, set := range sets {
		for _, src := range set.Sources() {
			ret.sources = append(ret.sources, &src)
		}
		for name, fn := range set.Funcs() {
			ret.funcs[name] = fn
		}
	}
	return ret
}

// Set appends a template source, identified by its file path and name,
// to the set. A copy of the set is returned.
func (s *Set) Set(name, file string) *Set {
	ret := s.clone()
	ret.sources = append(ret.sources, &Source{Name: name, FilePath: file})
	return ret
}

// AddFunc appends a named function to the set's function map. It is legal to
// overwrite elements of the map. A copy of the set is returned.
func (s *Set) AddFunc(name string, fn interface{}) *Set {
	ret := s.clone()
	ret.funcs[name] = fn
	return ret
}

// AddFuncs appends the functions of the passed-in function map to
// the set's function map. It is legal to overwrite elements of the map.
// A copy of the set is returned.
func (s *Set) AddFuncs(funcMap map[string]interface{}) *Set {
	ret := s.clone()
	for name, fn := range funcMap {
		ret.funcs[name] = fn
	}
	return ret
}

// Funcs returns a copy of the set's function map.
func (s *Set) Funcs() map[string]interface{} {
	ret := make(map[string]interface{}, len(s.funcs))
	for name, fn := range s.funcs {
		ret[name] = fn
	}
	return ret
}

// Sources returns a copy of the set's template sources.
func (s *Set) Sources() []Source {
	ret := make([]Source, len(s.sources))
	for i, src := range s.sources {
		ret[i] = *src
	}
	return ret
}

// View creates a new view from the set's template sources.
// If an error occurs, the returned view is nil.
func (s *Set) View() (*View, error) {
	view := &View{set: s}
	err := s.create(view)
	return view, err
}

// ViewMust creates a new view from the set's template sources.
// It panics if an error occurs.
func (s *Set) ViewMust() *View {
	view, err := s.View()
	if err != nil {
		panic(err)
	}
	return view
}

func (s *Set) clone() *Set {
	copy := *s
	return &copy
}

func (s *Set) create(v *View) error {
	if v.template != nil && !s.loader.conf.AutoReload {
		return nil
	}

	parsedTmpl, err := s.loader.buildTemplate(s)
	if err != nil {
		return err
	}
	v.template = parsedTmpl

	return nil
}

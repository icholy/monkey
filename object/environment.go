package object

type Env struct {
	parent *Env
	store  map[string]Object
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		store:  map[string]Object{},
	}
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return obj, ok
}

func (e *Env) Update(name string, val Object) bool {
	_, ok := e.store[name]
	if !ok && e.parent != nil {
		return e.parent.Update(name, val)
	}
	return ok
}

func (e *Env) Set(name string, val Object) {
	e.store[name] = val
}

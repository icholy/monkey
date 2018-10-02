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
	if obj, ok := e.store[name]; ok {
		return obj, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, true
}

func (e *Env) Set(name string, val Object) {
	e.store[name] = val
}

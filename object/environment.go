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
	if name == "locals" {
		return e.Locals(), true
	}
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
	e.store[name] = val
	return ok
}

func (e *Env) Set(name string, val Object) {
	e.store[name] = val
}

func (e *Env) Locals() Object {
	hash := NewHash()
	for k, v := range e.store {
		hash.Set(&String{Value: k}, v)
	}
	return hash
}

package object

type Env struct {
	store map[string]Object
}

func NewEnv() *Env {
	return &Env{
		store: map[string]Object{},
	}
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Env) Set(name string, val Object) {
	e.store[name] = val
}

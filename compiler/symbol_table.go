package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
	FreeScope    SymbolScope = "FREE"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer *SymbolTable
	store map[string]Symbol
	Count int
	Free  []Symbol
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Outer: outer,
		store: map[string]Symbol{},
	}
}

func (st *SymbolTable) DefineBuiltin(name string, index int) Symbol {
	s := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	st.store[name] = s
	return s
}

func (st *SymbolTable) Define(name string) Symbol {
	s := Symbol{Name: name, Index: st.Count, Scope: LocalScope}
	if st.Outer == nil {
		s.Scope = GlobalScope
	}
	st.Count++
	st.store[name] = s
	return s
}

func (st *SymbolTable) defineFree(original Symbol) Symbol {
	st.Free = append(st.Free, original)
	s := Symbol{
		Name:  original.Name,
		Scope: FreeScope,
		Index: len(st.Free) - 1,
	}
	st.store[s.Name] = s
	return s
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	s, ok := st.store[name]
	if !ok && st.Outer != nil {
		s, ok = st.Outer.Resolve(name)
		if !ok {
			return s, ok
		}
		if s.Scope == GlobalScope || s.Scope == BuiltinScope {
			return s, ok
		}
		return st.defineFree(s), true
	}
	return s, ok
}

package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	outer *SymbolTable
	store map[string]Symbol
	count int
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		outer: outer,
		store: map[string]Symbol{},
	}
}

func (st *SymbolTable) Define(name string) Symbol {
	s := Symbol{Name: name, Index: st.count, Scope: LocalScope}
	if st.outer == nil {
		s.Scope = GlobalScope
	}
	st.count++
	st.store[name] = s
	return s
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	s, ok := st.store[name]
	if !ok && st.outer != nil {
		return st.outer.Resolve(name)
	}
	return s, ok
}

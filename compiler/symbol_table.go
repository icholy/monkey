package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store map[string]Symbol
	count int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: map[string]Symbol{},
	}
}

func (st *SymbolTable) Define(name string) Symbol {
	s := Symbol{Name: name, Scope: GlobalScope, Index: st.count}
	st.count++
	st.store[name] = s
	return s
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	s, ok := st.store[name]
	return s, ok
}

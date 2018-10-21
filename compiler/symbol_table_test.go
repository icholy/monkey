package compiler

import (
	"testing"

	"gotest.tools/assert"
)

func TestDefineResolve(t *testing.T) {

	symbols := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "b", Scope: GlobalScope, Index: 2},
	}
	global := NewSymbolTable(nil)

	for _, s := range symbols {
		assert.Equal(t, global.Define(s.Name), s)

		actual, ok := global.Resolve(s.Name)
		assert.Assert(t, ok)
		assert.Equal(t, actual, s)
	}

}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable(nil)
	global.Define("a")
	global.Define("b")

	local := NewSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, symbol := range expected {
		actual, ok := local.Resolve(symbol.Name)
		assert.Assert(t, ok)
		assert.Equal(t, symbol, actual)
	}
}

func TestDefine(t *testing.T) {
	global := NewSymbolTable(nil)
	assert.Equal(t, global.Define("a"), Symbol{Name: "a", Scope: GlobalScope, Index: 0})
	assert.Equal(t, global.Define("b"), Symbol{Name: "b", Scope: GlobalScope, Index: 1})

	first := NewSymbolTable(global)
	assert.Equal(t, first.Define("c"), Symbol{Name: "c", Scope: LocalScope, Index: 0})
	assert.Equal(t, first.Define("d"), Symbol{Name: "d", Scope: LocalScope, Index: 1})

	second := NewSymbolTable(first)
	assert.Equal(t, second.Define("e"), Symbol{Name: "e", Scope: LocalScope, Index: 0})
	assert.Equal(t, second.Define("f"), Symbol{Name: "f", Scope: LocalScope, Index: 1})
}

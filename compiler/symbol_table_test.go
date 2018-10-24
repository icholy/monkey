package compiler

import (
	"testing"

	is "gotest.tools/assert/cmp"

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

func TestResolveFree(t *testing.T) {
	global := NewSymbolTable(nil)
	global.Define("a")
	global.Define("b")

	first := NewSymbolTable(global)
	first.Define("c")
	first.Define("d")

	second := NewSymbolTable(first)
	second.Define("e")
	second.Define("f")

	tests := []struct {
		table   *SymbolTable
		symbols []Symbol
		free    []Symbol
	}{
		{
			table: first,
			symbols: []Symbol{
				Symbol{Name: "a", Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Scope: GlobalScope, Index: 1},
				Symbol{Name: "c", Scope: LocalScope, Index: 0},
				Symbol{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			table: second,
			symbols: []Symbol{
				Symbol{Name: "a", Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Scope: GlobalScope, Index: 1},
				Symbol{Name: "c", Scope: FreeScope, Index: 0},
				Symbol{Name: "d", Scope: FreeScope, Index: 1},
				Symbol{Name: "e", Scope: LocalScope, Index: 0},
				Symbol{Name: "f", Scope: LocalScope, Index: 1},
			},
			free: []Symbol{
				Symbol{Name: "c", Scope: LocalScope, Index: 0},
				Symbol{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, expected := range tt.symbols {
			actual, ok := tt.table.Resolve(expected.Name)
			assert.Assert(t, ok, "symbol not found: %s", expected.Name)
			assert.Equal(t, actual, expected)
		}
		assert.Assert(t, is.Len(tt.table.Free, len(tt.free)))
		for i, expected := range tt.free {
			assert.Equal(t, expected, tt.table.Free[i])
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := NewSymbolTable(nil)
	global.Define("a")

	first := NewSymbolTable(global)
	first.Define("c")

	second := NewSymbolTable(first)
	second.Define("e")
	second.Define("f")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "c", Scope: FreeScope, Index: 0},
		{Name: "e", Scope: LocalScope, Index: 0},
		{Name: "f", Scope: LocalScope, Index: 1},
	}
	for _, expected := range expected {
		actual, ok := second.Resolve(expected.Name)
		assert.Assert(t, ok, "symbol not found: %s", expected.Name)
		assert.Equal(t, actual, expected)
	}

	unexpected := []string{"b", "d"}
	for _, unexpected := range unexpected {
		_, ok := second.Resolve(unexpected)
		assert.Assert(t, !ok, "unexpected symbol found: %s", unexpected)
	}
}

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
	global := NewSymbolTable()

	for _, s := range symbols {
		assert.Equal(t, global.Define(s.Name), s)

		actual, ok := global.Resolve(s.Name)
		assert.Assert(t, ok)
		assert.Equal(t, actual, s)
	}

}

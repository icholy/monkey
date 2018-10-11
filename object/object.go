package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/icholy/monkey/ast"
)

type ObjectType string

const (
	INTEGER  = "INTEGER"
	NULL     = "NULL"
	BOOLEAN  = "BOOLEAN"
	RETURN   = "RETURN"
	FUNCTION = "FUNCTION"
	STRING   = "STRING"
	BUILTIN  = "BUILTIN"
	ARRAY    = "ARRAY"
	HASH     = "HASH"
)

var types = map[string]ObjectType{
	"integer":  INTEGER,
	"boolean":  BOOLEAN,
	"string":   STRING,
	"array":    ARRAY,
	"hash":     HASH,
	"function": FUNCTION,
}

func LookupType(name string) (ObjectType, bool) {
	t, ok := types[name]
	return t, ok
}

type Object interface {
	Type() ObjectType
	Inspect() string
	KeyValue() KeyValue
}

type TypedObject struct {
	Object
	ObjectType ObjectType
}

func (o *TypedObject) Set(val Object) error {
	if val.Type() != o.ObjectType {
		return fmt.Errorf("wrong type: expected %s, got %s", o.ObjectType, val.Type())
	}
	return nil
}

type BuiltinFunc func(...Object) (Object, error)

type Builtin struct {
	Fn BuiltinFunc
}

func (b *Builtin) KeyValue() KeyValue { return b.Fn }
func (b *Builtin) Inspect() string    { return "builtin function" }
func (b *Builtin) Type() ObjectType   { return BUILTIN }

type Integer struct {
	Value int64
}

func (i *Integer) KeyValue() KeyValue { return i.Value }
func (i *Integer) Inspect() string    { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType   { return INTEGER }

type Boolean struct {
	Value bool
}

func (b *Boolean) KeyValue() KeyValue { return b.Value }
func (b *Boolean) Inspect() string    { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() ObjectType   { return BOOLEAN }

type String struct {
	Value string
}

func (s *String) At(i int) (Object, error) {
	if i < 0 || i >= len(s.Value) {
		return nil, fmt.Errorf("%d out of range", i)
	}
	return &String{Value: string(s.Value[i])}, nil
}

func (s *String) KeyValue() KeyValue { return s.Value }
func (s *String) Inspect() string    { return fmt.Sprintf("%q", s.Value) }
func (s *String) Type() ObjectType   { return STRING }

type Null struct{}

func (n *Null) KeyValue() KeyValue { return nil }
func (n *Null) Inspect() string    { return "null" }
func (n *Null) Type() ObjectType   { return NULL }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) KeyValue() KeyValue { return r.Value.KeyValue() }
func (r *ReturnValue) Inspect() string    { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType   { return RETURN }

func UnwrapReturn(obj Object) Object {
	if ret, ok := obj.(*ReturnValue); ok {
		return ret.Value
	}
	return obj
}

type Function struct {
	Parameters []*ast.Parameter
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) KeyValue() KeyValue { return f }

func (f *Function) Type() ObjectType { return FUNCTION }

func (f *Function) Inspect() string {
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.Name.Value)
	}
	return fmt.Sprintf("fn(%s)", strings.Join(params, ", "))
}

type Array struct {
	Elements []Object
}

func (a *Array) InRange(i int) bool {
	return i >= 0 && i < len(a.Elements)
}

func (a *Array) At(i int) (Object, error) {
	if !a.InRange(i) {
		return nil, fmt.Errorf("%d not in range", i)
	}
	return a.Elements[i], nil
}

func (a *Array) SetAt(i int, v Object) error {
	if !a.InRange(i) {
		return fmt.Errorf("%d not in range", i)
	}
	a.Elements[i] = v
	return nil
}

func (a *Array) KeyValue() KeyValue { return a }
func (Array) Type() ObjectType      { return ARRAY }
func (a *Array) Inspect() string {
	var vals []string
	for _, e := range a.Elements {
		vals = append(vals, e.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(vals, ", "))
}

type KeyValue interface{}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	pairs map[KeyValue]*HashPair
}

func NewHash() *Hash {
	return &Hash{
		pairs: map[KeyValue]*HashPair{},
	}
}

func (h *Hash) Set(key, value Object) {
	h.pairs[key.KeyValue()] = &HashPair{
		Key:   key,
		Value: value,
	}
}

func (h *Hash) SetPairs(pairs ...*HashPair) {
	for _, p := range pairs {
		h.Set(p.Key, p.Value)
	}
}

func (h *Hash) Get(key Object) (Object, bool) {
	p, ok := h.pairs[key.KeyValue()]
	if !ok {
		return nil, false
	}
	return p.Value, true
}

func (h *Hash) Delete(key Object) {
	delete(h.pairs, key.KeyValue())
}

func (h *Hash) Len() int {
	return len(h.pairs)
}

func (h *Hash) Pairs() []*HashPair {
	var pairs []*HashPair
	for _, p := range h.pairs {
		pairs = append(pairs, p)
	}
	return pairs
}

func (h *Hash) KeyValue() KeyValue { return h }
func (Hash) Type() ObjectType      { return HASH }
func (h *Hash) Inspect() string {
	if h.Len() == 0 {
		return "{}"
	}
	var pairs []string
	for _, p := range h.pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", p.Key.Inspect(), p.Value.Inspect()))
	}
	return fmt.Sprintf("{ %s }", strings.Join(pairs, ", "))
}

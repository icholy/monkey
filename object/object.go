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
	ERROR    = "ERROR"
	FUNCTION = "FUNCTION"
	STRING   = "STRING"
	BUILTIN  = "BUILTIN"
	ARRAY    = "ARRAY"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type BuiltinFunc func(...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunc
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return fmt.Sprintf("%v", s.Value) }
func (s *String) Type() ObjectType { return STRING }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType { return RETURN }

type Error struct {
	Message string
}

func (e *Error) Error() string    { return e.Message }
func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return fmt.Sprintf("ERROR: %s", e.Message) }

func Errorf(format string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
	}
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) Type() ObjectType { return FUNCTION }

func (f *Function) Inspect() string {
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.Value)
	}
	return fmt.Sprintf("fn(%s) %s", strings.Join(params, ", "), f.Body)
}

type Array struct {
	Elements []Object
}

func (Array) Type() ObjectType { return ARRAY }
func (a *Array) Inspect() string {
	var vals []string
	for _, e := range a.Elements {
		vals = append(vals, e.Inspect())
	}
	return strings.Join(vals, ", ")
}

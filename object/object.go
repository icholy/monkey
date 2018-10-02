package object

import (
	"fmt"
	"strconv"
)

type ObjectType string

const (
	INTEGER = "INTEGER"
	NULL    = "NULL"
	BOOLEAN = "BOOLEAN"
	RETURN  = "RETURN"
	ERROR   = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

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

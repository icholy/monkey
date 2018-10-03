package evaluator

import (
	"github.com/icholy/monkey/object"
)

var builtins = map[string]object.Object{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("len: wrong number of arguments")
			}
			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(obj.Value))}
			default:
				return object.Errorf("len: invalid argument type %s", args[0].Type())
			}
		},
	},
}

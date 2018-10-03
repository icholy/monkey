package evaluator

import (
	"fmt"

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
			case *object.Array:
				return &object.Integer{Value: int64(len(obj.Elements))}
			default:
				return object.Errorf("len: invalid argument type %s", args[0].Type())
			}
		},
	},
	"append": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return object.Errorf("append: at least one argument required")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return object.Errorf("append: expected array, got %s", args[0].Type())
			}
			arr.Elements = append(arr.Elements, args[1:]...)
			return arr
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			var values []interface{}
			for _, a := range args {
				values = append(values, a.Inspect())
			}
			fmt.Println(values...)
			return &object.Null{}
		},
	},
}

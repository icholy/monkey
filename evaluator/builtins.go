package evaluator

import (
	"fmt"
	"io/ioutil"

	"github.com/icholy/monkey/object"
)

var builtins = map[string]object.Object{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("len: wrong number of arguments")
			}
			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(obj.Value))}, nil
			case *object.Array:
				return &object.Integer{Value: int64(len(obj.Elements))}, nil
			default:
				return nil, fmt.Errorf("len: invalid argument type %s", args[0].Type())
			}
		},
	},
	"append": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("append: at least one argument required")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return nil, fmt.Errorf("append: expected array, got %s", args[0].Type())
			}
			arr.Elements = append(arr.Elements, args[1:]...)
			return arr, nil
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("first: expecting one array argument")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return nil, fmt.Errorf("first: expecting one array argument")
			}
			if len(arr.Elements) == 0 {
				return nil, fmt.Errorf("first: cannot get first element of empty array")
			}
			return arr.Elements[0], nil
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("first: expecting one array argument")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return nil, fmt.Errorf("first: expecting one array argument")
			}
			if len(arr.Elements) == 0 {
				return &object.Array{}, nil
			}
			return &object.Array{
				Elements: arr.Elements[1:],
			}, nil
		},
	},
	"print": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			var values []interface{}
			for _, a := range args {
				values = append(values, a.Inspect())
			}
			fmt.Println(values...)
			return &object.Null{}, nil
		},
	},
	"read": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("read: requires one argument")
			}
			name, ok := args[0].(*object.String)
			if !ok {
				return nil, fmt.Errorf("read: expected a string, got %s", args[0].Type())
			}
			data, err := ioutil.ReadFile(name.Value)
			if err != nil {
				return nil, fmt.Errorf("read: %v", err)
			}
			return &object.String{
				Value: string(data),
			}, nil
		},
	},
	"keys": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("keys: requires one argument")
			}
			hash, ok := args[0].(*object.Hash)
			if !ok {
				return nil, fmt.Errorf("keys: expected a hash, got %s", args[0].Type())
			}
			arr := &object.Array{}
			for _, p := range hash.Pairs() {
				arr.Elements = append(arr.Elements, p.Key)
			}
			return arr, nil
		},
	},
	"values": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("keys: requires one argument")
			}
			hash, ok := args[0].(*object.Hash)
			if !ok {
				return nil, fmt.Errorf("keys: expected a hash, got %s", args[0].Type())
			}
			arr := &object.Array{}
			for _, p := range hash.Pairs() {
				arr.Elements = append(arr.Elements, p.Value)
			}
			return arr, nil
		},
	},
	"str": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("str: requires one argument")
			}
			if s, ok := args[0].(*object.String); ok {
				return s, nil
			}
			return &object.String{
				Value: args[0].Inspect(),
			}, nil
		},
	},
	"type": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("type: requires one argument")
			}
			return &object.String{
				Value: string(args[0].Type()),
			}, nil
		},
	},
}

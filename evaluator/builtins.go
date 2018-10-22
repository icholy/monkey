package evaluator

import (
	"fmt"
	"io/ioutil"
	"reflect"

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
			case *object.Hash:
				return &object.Integer{Value: int64(obj.Len())}, nil
			default:
				return nil, fmt.Errorf("len: invalid argument type %s", args[0].Type())
			}
		},
	},
	"delete": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("delete: expecting two arguments")
			}
			hash, ok := args[0].(*object.Hash)
			if !ok {
				return nil, fmt.Errorf("delete: expecting first parameter to be hash")
			}
			hash.Delete(args[1])
			return NULL, nil
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
				values = append(values, a.Inspect(0))
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
				Value: args[0].Inspect(0),
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

func WrapFunc(fn interface{}) object.BuiltinFunc {
	var (
		objectType = reflect.TypeOf((*object.Object)(nil)).Elem()
		errorType  = reflect.TypeOf((*error)(nil)).Elem()
		fnType     = reflect.TypeOf(fn)
		fnValue    = reflect.ValueOf(fn)
	)

	// make sure it's a function a usable function
	if fnType.Kind() != reflect.Func {
		panic("builtin: must be a function")
	}
	if fnType.IsVariadic() {
		panic("builtin: variadic functions are not supported")
	}

	// check the return types
	if fnType.NumOut() != 2 {
		panic("builtin: expected 2 return values")
	}
	if fnType.Out(0) != objectType {
		panic("builtin: return value 1 should be object.Object")
	}
	if fnType.Out(1) != errorType {
		panic("builtin: return value 2 should be error")
	}

	// check the parameters
	var params []reflect.Type
	for i := 0; i < fnType.NumIn(); i++ {
		if !fnType.In(i).Implements(objectType) {
			panic("builtin: param doesn't implement object.Object")
		}
		params = append(params, fnType.In(i))
	}

	return func(args ...object.Object) (object.Object, error) {
		// check the parameters
		if len(args) != len(params) {
			return nil, fmt.Errorf("wrong number of arguments")
		}

		in := make([]reflect.Value, len(params))
		for i, arg := range args {
			value := reflect.ValueOf(arg)
			if !value.Type().AssignableTo(params[i]) {
				return nil, fmt.Errorf("invalid argument: %d %s", i, arg.Inspect(0))
			}
			in[i] = value
		}
		out := fnValue.Call(in)
		return out[0].Interface().(object.Object), out[1].Interface().(error)
	}
}

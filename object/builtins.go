package object

import (
	"fmt"
	"io/ioutil"
	"reflect"
)

type BuiltinFunc func(...Object) (Object, error)

type Builtin struct {
	Name string
	Fn   BuiltinFunc
}

func (b *Builtin) KeyValue() KeyValue       { return b.Fn }
func (b *Builtin) Inspect(depth int) string { return "<builtin function>" }
func (b *Builtin) Type() ObjectType         { return BUILTIN }

var Builtins = []*Builtin{
	&Builtin{
		Name: "len",
		Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("len: wrong number of arguments")
			}
			switch obj := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(obj.Value))}, nil
			case *Array:
				return &Integer{Value: int64(len(obj.Elements))}, nil
			case *Hash:
				return &Integer{Value: int64(obj.Len())}, nil
			default:
				return nil, fmt.Errorf("len: invalid argument type %s", args[0].Type())
			}
		},
	},
	&Builtin{
		Name: "delete",
		Fn: MakeBuiltinFunc(func(hash *Hash, key Object) (Object, error) {
			hash.Delete(key)
			return nil, nil
		}),
	},
	&Builtin{
		Name: "append",
		Fn: func(args ...Object) (Object, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("append: at least one argument required")
			}
			arr, ok := args[0].(*Array)
			if !ok {
				return nil, fmt.Errorf("append: expected array, got %s", args[0].Type())
			}
			arr.Elements = append(arr.Elements, args[1:]...)
			return arr, nil
		},
	},
	&Builtin{
		Name: "first",
		Fn: MakeBuiltinFunc(func(arr *Array) (Object, error) {
			if len(arr.Elements) == 0 {
				return nil, fmt.Errorf("first: cannot get first element of empty array")
			}
			return arr.Elements[0], nil
		}),
	},
	&Builtin{
		Name: "rest",
		Fn: MakeBuiltinFunc(func(arr *Array) (Object, error) {
			if len(arr.Elements) == 0 {
				return &Array{}, nil
			}
			return &Array{
				Elements: arr.Elements[1:],
			}, nil
		}),
	},
	&Builtin{
		Name: "print",
		Fn: func(args ...Object) (Object, error) {
			var values []interface{}
			for _, a := range args {
				values = append(values, a.Inspect(0))
			}
			fmt.Println(values...)
			return nil, nil
		},
	},
	&Builtin{
		Name: "read",
		Fn: MakeBuiltinFunc(func(name *String) (Object, error) {
			data, err := ioutil.ReadFile(name.Value)
			if err != nil {
				return nil, fmt.Errorf("read: %v", err)
			}
			return &String{
				Value: string(data),
			}, nil
		}),
	},
	&Builtin{
		Name: "keys",
		Fn: MakeBuiltinFunc(func(hash *Hash) (Object, error) {
			arr := &Array{}
			for _, p := range hash.Pairs() {
				arr.Elements = append(arr.Elements, p.Key)
			}
			return arr, nil
		}),
	},
	&Builtin{
		Name: "values",
		Fn: MakeBuiltinFunc(func(hash *Hash) (Object, error) {
			arr := &Array{}
			for _, p := range hash.Pairs() {
				arr.Elements = append(arr.Elements, p.Value)
			}
			return arr, nil
		}),
	},
	&Builtin{
		Name: "str",
		Fn: MakeBuiltinFunc(func(v Object) (Object, error) {
			if v.Type() == STRING {
				return v, nil
			}
			return &String{Value: v.Inspect(0)}, nil
		}),
	},
	&Builtin{
		Name: "type",
		Fn: MakeBuiltinFunc(func(v Object) (Object, error) {
			return &String{Value: string(v.Type())}, nil
		}),
	},
}

func BuiltinMap() map[string]*Builtin {
	m := map[string]*Builtin{}
	for _, b := range Builtins {
		m[b.Name] = b
	}
	return m
}

func FindBuiltin(name string) int {
	for i, b := range Builtins {
		if b.Name == name {
			return i
		}
	}
	return -1
}

func MakeBuiltinFunc(fn interface{}) BuiltinFunc {
	var (
		objectType = reflect.TypeOf((*Object)(nil)).Elem()
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
		panic("builtin: return value 1 should be Object")
	}
	if fnType.Out(1) != errorType {
		panic("builtin: return value 2 should be error")
	}

	// check the parameters
	var params []reflect.Type
	for i := 0; i < fnType.NumIn(); i++ {
		if !fnType.In(i).Implements(objectType) {
			panic("builtin: param doesn't implement Object")
		}
		params = append(params, fnType.In(i))
	}

	return func(args ...Object) (Object, error) {
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

		var result Object
		if v := out[0]; !v.IsNil() {
			result = v.Interface().(Object)
		}
		var err error
		if v := out[1]; !v.IsNil() {
			err = v.Interface().(error)
		}
		return result, err
	}
}

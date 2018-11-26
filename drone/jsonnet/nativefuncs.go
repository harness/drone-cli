package jsonnet

import (
	"bytes"
	"encoding/json"
	"io"

	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// RegisterNativeFuncs adds kubecfg's native jsonnet functions to provided VM
func RegisterNativeFuncs(vm *jsonnet.VM) {
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "parseJson",
		Params: []jsonnetAst.Identifier{"json"},
		Func: func(args []interface{}) (res interface{}, err error) {
			data := []byte(args[0].(string))
			err = json.Unmarshal(data, &res)
			return
		},
	})

	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "parseYaml",
		Params: []jsonnetAst.Identifier{"yaml"},
		Func: func(args []interface{}) (res interface{}, err error) {
			ret := []interface{}{}
			data := []byte(args[0].(string))
			d := yaml.NewYAMLToJSONDecoder(bytes.NewReader(data))
			for {
				var doc interface{}
				if err := d.Decode(&doc); err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}
				ret = append(ret, doc)
			}
			return ret, nil
		},
	})
}

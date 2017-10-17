package transistor

import (
	"reflect"
)

type Creator func() Plugin

var PluginRegistry = map[string]Creator{}
var ApiRegistry = make(map[string]interface{})

func RegisterPlugin(name string, creator Creator) {
	PluginRegistry[name] = creator
}

func RegisterApi(i interface{}) {
	ApiRegistry[reflect.TypeOf(i).String()] = i
}

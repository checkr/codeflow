package plugins

import "github.com/codeamp/transistor"

func init() {
	transistor.RegisterApi(Hello{})
}

type Hello struct {
	Action  string
	Message string
}

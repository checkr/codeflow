package schema

import (
	"io/ioutil"
	"path"
	"runtime"
)

func Schema() (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	schema, err := ioutil.ReadFile(path.Join(path.Dir(filename), "schema.graphql"))
	if err != nil {
		return "", err
	}

	return string(schema), nil
}

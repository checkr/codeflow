package json

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

type singularExample struct {
  A int `json:"a" role:"public,private"`
  B int `json:"b" role:"private"`
  C int `json:"c"`
}

type rolesExample struct {
  A int `json:"a" roles:"public,private"`
  B int `json:"b" roles:"private"`
  C int `json:"c"`
  D singularExample `json:"d"`
}

func TestBasicRoles(t *testing.T) {
  var data []byte
  var err error
  
  ex1 := rolesExample{1, 2, 3, singularExample{4, 5, 6}}
  
  data, err = MarshalRole("public", ex1)
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, `{"a":1,"c":3,"d":{"a":4,"c":6}}`, string(data))
  }
  
  data, err = MarshalRole("private", ex1)
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, `{"a":1,"b":2,"c":3,"d":{"a":4,"b":5,"c":6}}`, string(data))
  }
  
}


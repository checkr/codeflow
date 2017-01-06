# Go JSON

This is a drop-in replacement JSON marshaler/encoder which supports the conditional exclusion of properties.

You can use this package instead of the standard library JSON package in cases where you want to control which properties are included in the marshaled result. This package is a modification of the standard library JSON encoder and should otherwise function exactly the same.

    import (
      "github.com/bww/json"
    )
    
    type Example struct {
      A int `json:"a" roles:"public,private"`
      B int `json:"b" roles:"private"`
      C int `json:"c"`
    }
    
    func enc() {
      var data []byte
      
      ex := Example{
        1, 2, 3,
      }
      
      // only properties with "public" or undefined roles are included
      data, _ = json.MarshalRole("public", ex)
      fmt.Println(string(data)) // {"a":1,"C",3}
      
      // only properties with "private" or undefined roles are included
      data, _ = json.MarshalRole("private", ex)
      fmt.Println(string(data)) // {"a":1,"b":2,"c",3}
      
      // for compatibility with encoding/json, roles are ignored
      data, _ = json.Marshal(ex)
      fmt.Println(string(data)) // {"a":1,"b":2,"c",3}
      
    }

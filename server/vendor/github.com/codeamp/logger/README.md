# Simple logrus wrapper

### Example Usage
```
package main

import (
	log "github.com/codeamp/logger"
)

func main() {
  log.InfoWithFields("hello", log.Fields{
    "world": "earth",
  })

  log.Println("Hello World")
  
  log.Debug("Hello World")
}
```

[GoDoc](https://godoc.org/github.com/codeamp/logger)

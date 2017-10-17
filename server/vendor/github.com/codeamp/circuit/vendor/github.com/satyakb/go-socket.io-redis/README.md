# go-socket.io-redis

By running go-socket.io with this adapter, you can run multiple socket.io 
instances in different processes or servers that can all broadcast and emit 
events to and from each other.

## How to use

Install the package using:

```bash
go get "github.com/satyakb/go-socket.io-redis"
```

Usage:

```go
import (
    "log"
    "net/http"
    "github.com/googollee/go-socket.io"
    "github.com/satyakb/go-socket.io-redis"
)

func main() {
    server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }

    opts := make(map[string]string)
    server.SetAdaptor(redis.Redis(opts))

    server.On("connection", func(so socketio.Socket) {
        log.Println("on connection")
        so.Join("chat")
        so.On("chat message", func(msg string) {
            log.Println("emit:", so.Emit("chat message", msg))
            so.BroadcastTo("chat", "chat message", msg)
        })
        so.On("disconnection", func() {
            log.Println("on disconnect")
        })
    })
    server.On("error", func(so socketio.Socket, err error) {
        log.Println("error:", err)
    })

    http.Handle("/socket.io/", server)
    http.Handle("/", http.FileServer(http.Dir("./asset")))
    log.Println("Serving at localhost:5000...")
    log.Fatal(http.ListenAndServe(":5000", nil))
}
```

**Note:** The package is named `redis` for use in your code

## API

### Redis(opts map[string]string)

The following options are allowed:

- `host`: host to connect to redis on (`"localhost"`)
- `port`: port to connect to redis on (`"6379"`)
- `prefix`: the prefix of the key to pub/sub events on (`"socket.io"`)

## References

Code and README based off of:
- https://github.com/Automattic/socket.io-redis


package redis

import (
    "log"
    "strings"
    "github.com/googollee/go-socket.io"
    "github.com/garyburd/redigo/redis"
    "github.com/nu7hatch/gouuid"
    "github.com/satyakb/go-socket.io-redis/cmap_string_cmap"
    "github.com/satyakb/go-socket.io-redis/cmap_string_socket"
    // "github.com/vmihailenco/msgpack"  // screwed up types after decoding
    "encoding/json"
)

type broadcast struct {
  host string
  port string
  pub redis.PubSubConn
  sub redis.PubSubConn
  prefix string
  uid string
  key string
  remote bool
  rooms cmap_string_cmap.ConcurrentMap
}

//
// opts: {
//   "host": "127.0.0.1",
//   "port": "6379"
//   "prefix": "socket.io"
// }
func Redis(opts map[string]string) socketio.BroadcastAdaptor {
  b := broadcast {
    rooms: cmap_string_cmap.New(),
  }

  var ok bool
  b.host, ok = opts["host"]
  if !ok {
    b.host = "127.0.0.1"
  }
  b.port, ok = opts["port"]
  if !ok {
    b.port = "6379"
  }
  b.prefix, ok = opts["prefix"]
  if !ok {
    b.prefix = "socket.io"
  }

  pub, err := redis.Dial("tcp", b.host + ":" + b.port)
  if err != nil {
      panic(err)
  }
  sub, err := redis.Dial("tcp", b.host + ":" + b.port)
  if err != nil {
      panic(err)
  }

  b.pub = redis.PubSubConn{Conn: pub}
  b.sub = redis.PubSubConn{Conn: sub}

  uid, err := uuid.NewV4();
  if err != nil {
    log.Println("error generating uid:", err)
    return nil
  }
  b.uid = uid.String()
  b.key = b.prefix + "#" + b.uid

  b.remote = false

  b.sub.PSubscribe(b.prefix + "#*")

  // This goroutine receives and prints pushed notifications from the server.
  // The goroutine exits when there is an error.
  go func() {
      for {
          switch n := b.sub.Receive().(type) {
          case redis.Message:
              log.Printf("Message: %s %s\n", n.Channel, n.Data)
          case redis.PMessage:
              b.onmessage(n.Channel, n.Data)
              log.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
          case redis.Subscription:
              log.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
              if n.Count == 0 {
                  return
              }
          case error:
              log.Printf("error: %v\n", n)
              return
          }
      }
  }()

  return b
}

func (b broadcast) onmessage(channel string, data []byte) error {
  pieces := strings.Split(channel, "#");
  uid := pieces[len(pieces) - 1]
  if b.uid == uid {
    log.Println("ignore same uid")
    return nil
  }

  var out map[string][]interface{}
  err := json.Unmarshal(data, &out)
  if err != nil {
    log.Println("error decoding data")
    return nil
  }

  args := out["args"]
  opts := out["opts"]
  ignore, ok := opts[0].(socketio.Socket)
  if !ok {
    log.Println("ignore is not a socket")
    ignore = nil
  }
  room, ok := opts[1].(string)
  if !ok {
    log.Println("room is not a string")
    room = ""
  }
  message, ok := opts[2].(string)
  if !ok {
    log.Println("message is not a string")
    message = ""
  }

  b.remote = true;
  b.Send(ignore, room, message, args...)
  return nil
}

func (b broadcast) Join(room string, socket socketio.Socket) error {
  sockets, ok := b.rooms.Get(room)
  if !ok {
    sockets = cmap_string_socket.New()
  }
  sockets.Set(socket.Id(), socket)
  b.rooms.Set(room, sockets)
  return nil
}

func (b broadcast) Leave(room string, socket socketio.Socket) error {
  sockets, ok := b.rooms.Get(room)
  if !ok {
    return nil
  }
  sockets.Remove(socket.Id())
  if sockets.IsEmpty() {
    b.rooms.Remove(room)
    return nil
  }
  b.rooms.Set(room, sockets)
  return nil
}

// Same as Broadcast
func (b broadcast) Send(ignore socketio.Socket, room, message string, args ...interface{}) error {
  sockets, ok := b.rooms.Get(room)
  if !ok {
    return nil
  }
  for item := range sockets.Iter() {
    id := item.Key
    s := item.Val
    if ignore != nil && ignore.Id() == id {
      continue
    }
    err := (s.Emit(message, args...))
    if err != nil {
      log.Println("error broadcasting:", err)
    }
  }

  opts := make([]interface{}, 3)
  opts[0] = ignore
  opts[1] = room
  opts[2] = message
  in := map[string][]interface{}{
    "args": args,
    "opts": opts,
  }

  buf, err := json.Marshal(in)
  _ = err

  if !b.remote {
    b.pub.Conn.Do("PUBLISH", b.key, buf)
  }
  b.remote = false
  return nil
}

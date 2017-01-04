package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/checkr/codeflow/server/agent"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Websockets struct {
	ServiceAddress string `mapstructure:"service_address"`

	events chan agent.Event

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func init() {
	agent.RegisterPlugin("websockets", func() agent.Plugin {
		return &Websockets{}
	})
}

func (x *Websockets) Description() string {
	return "Send events subscribed clients"
}

func (x *Websockets) SampleConfig() string {
	return ` `
}

func (x *Websockets) Listen() {
	go x.run()

	r := mux.NewRouter()
	r.HandleFunc("/", x.serveWs)

	err := http.ListenAndServe(fmt.Sprintf("%s", x.ServiceAddress), r)
	if err != nil {
		log.Printf("Error starting server: %v", err)
	}
}

// serveWs handles websocket requests from the peer.
func (x *Websockets) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{ws: x, conn: conn, send: make(chan []byte, 256)}
	client.ws.register <- client
	go client.writePump()
	client.readPump()
}

func (x *Websockets) Start(e chan agent.Event) error {
	x.events = e
	x.broadcast = make(chan []byte)
	x.register = make(chan *Client)
	x.unregister = make(chan *Client)
	x.clients = make(map[*Client]bool)

	go x.Listen()
	log.Printf("Started the Websockets service on %s\n", x.ServiceAddress)

	return nil
}

func (x *Websockets) Stop() {
	log.Println("Stopping Websockets")
}

func (x *Websockets) Subscribe() []string {
	return []string{
		"plugins.WebsocketMsg",
	}
}

func (x *Websockets) Process(e agent.Event) error {
	log.Printf("Process Websockets event: %s", e.Name)
	if e.Name == "plugins.WebsocketMsg" {
		json, _ := json.Marshal(e.Payload)
		x.broadcast <- json
	}
	return nil
}

func (x *Websockets) run() {
	for {
		select {
		case client := <-x.register:
			x.clients[client] = true
		case client := <-x.unregister:
			if _, ok := x.clients[client]; ok {
				delete(x.clients, client)
				close(client.send)
			}
		case message := <-x.broadcast:
			for client := range x.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(x.clients, client)
				}
			}
		}
	}
}

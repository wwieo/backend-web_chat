package chat

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/utils"

	"github.com/gorilla/websocket"
)

type Chat struct {
	users    map[string]*User
	messages chan *Message
	join     chan *User
	leave    chan *User
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("%s %s%s %v\n", r.Method, r.Host, r.RequestURI, r.Proto)
		return r.Method == http.MethodGet
	},
}

func (c *Chat) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error on websocket connection:", err.Error())
	}

	keys := r.URL.Query()
	username := keys.Get("username")
	if strings.TrimSpace(username) == "" {
		username = fmt.Sprintf("anom-%d", utils.GetRandomI64())
	}

	c.join <- &User{
		UserName: username,
		Conn:     conn,
		Global:   c,
	}
}

func (c *Chat) Run() {
	for {
		select {
		case user := <-c.join:
			c.add(user)
		}
	}
}

func (c *Chat) add(user *User) {
	if _, ok := c.users[user.UserName]; !ok {
		c.users[user.UserName] = user
		log.Printf("Added user: %s, Total: %d\n", user.UserName, len(c.users))
	}
}

func Start(port string) {
	log.Printf("Chat listening on port: %s\n", port)
	c := &Chat{
		users:    make(map[string]*User),
		messages: make(chan *Message),
		join:     make(chan *User),
		leave:    make(chan *User),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(([]byte("Welcome to Go Webchat!")))
	})
	http.HandleFunc("/chat", c.Handler)

	go c.Run()

	log.Fatal(http.ListenAndServe(port, nil))
}

package main

import (
	"github.com/mxkacsa/eventbus"
	"log"
	"time"
)

type Message struct {
	Sender string
	Text   string
}

type Server struct {
	Eb *eventbus.EventBus[Message]
}

type Client struct {
	Name   string
	server *Server
}

func NewClient(name string, server *Server) *Client {
	return &Client{Name: name, server: server}
}

func (c *Client) receiveMessage(message Message) {
	if message.Sender == c.Name {
		return
	}

	log.Println(c.Name, "got message from ", message.Sender, " - ", message.Text)
}

func (c *Client) SendMessage(text string) {
	c.server.Eb.Publish(Message{
		Sender: c.Name,
		Text:   text,
	})
}

func (c *Client) Connect() error {
	log.Println(c.Name, "Connected")
	return c.server.Eb.SubscribeMethod(c, c.receiveMessage)
}

func (c *Client) Disconnect() error {
	log.Println(c.Name, "Disconnected")
	return c.server.Eb.UnsubscribeMethod(c, c.receiveMessage)
}

func main() {
	server := &Server{Eb: new(eventbus.EventBus[Message])}

	bob := NewClient("Bob", server)
	err := bob.Connect()
	if err != nil {
		return
	}

	alice := NewClient("Alice", server)
	err = alice.Connect()
	if err != nil {
		return
	}

	go func() {
		for i := 0; i < 3; i++ {
			if i%2 == 0 {
				bob.SendMessage("Hello!")
			} else {
				alice.SendMessage("Ola!")
			}

			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(5 * time.Second)

	err = bob.Disconnect()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	alice.SendMessage("Where are you?")

	time.Sleep(1 * time.Second)
}

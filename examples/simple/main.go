package main

import (
	"github.com/mxkacsa/eventbus"
	"log"
)

func main() {
	eb := new(eventbus.EventBus[string])
	err := eb.SubscribeFunc(callback)
	if err != nil {
		log.Fatal(err)
	}
	eb.Publish("hello")
	eb.Publish("world")

	err = eb.UnsubscribeFunc(callback)
	if err != nil {
		log.Fatal(err)
	}

	eb.Publish("!")
}

func callback(message string) {
	log.Println(message)
}

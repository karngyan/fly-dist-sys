package main

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("broadcast", handlerWithNode(broadcastHandler, n))
	n.Handle("read", handlerWithNode(readHandler, n))
	n.Handle("topology", handlerWithNode(topologyHandler, n))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

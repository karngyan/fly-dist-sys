package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var counter int64

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["id"] = getUniqueId(msg.Dest)
		body["type"] = "generate_ok"

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func getUniqueId(nid string) string {
	// timestamp - node id - counter
	ts := time.Now().UnixNano() / 1000000
	atomic.AddInt64(&counter, 1)
	return fmt.Sprintf("%d-%s-%d", ts, nid, counter)
}

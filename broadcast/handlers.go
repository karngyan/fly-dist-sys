package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type HandlerFunc func(msg maelstrom.Message, n *maelstrom.Node) error
type HandlerFuncWG func(msg maelstrom.Message, n *maelstrom.Node, wg *sync.WaitGroup) error

func handlerWithNode(h HandlerFunc, n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return h(msg, n)
	}
}

func handlerWithNodeWG(h HandlerFuncWG, n *maelstrom.Node, wg *sync.WaitGroup) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return h(msg, n, wg)
	}
}

var data = newStore()

type broadcastReq struct {
	Message int64  `json:"message"`
	MsgId   int64  `json:"msg_id"`
	Type    string `json:"type"`
}

func broadcastHandler(msg maelstrom.Message, n *maelstrom.Node) error {

	var body broadcastReq

	var resp struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if !data.getMByKey(body.Message) {
		go processBroadcastWithRetry(msg, n, body)
	}

	if body.MsgId == 0 {
		// as we know this is a Send call from one of nodes
		return nil
	}

	resp.Type = "broadcast_ok"
	return n.Reply(msg, resp)
}

func processBroadcastWithRetry(msg maelstrom.Message, n *maelstrom.Node, body broadcastReq) {
	data.addM(body.Message)
	neighbors := data.getNeighbors()

	unsent := map[string]bool{}
	for _, neighbor := range neighbors {
		unsent[neighbor] = false
	}

	for {
		if len(unsent) == 0 {
			break
		}

		for _, neighbor := range neighbors {
			dest := neighbor
			if dest == msg.Src {
				// no need to send it back to the sender
				delete(unsent, dest)
				continue
			}

			if err := n.RPC(dest, body, func(resp maelstrom.Message) error {
				delete(unsent, dest)
				return nil
			}); err != nil {
				log.Printf("Error sending message: %s", err)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func readHandler(msg maelstrom.Message, n *maelstrom.Node) error {
	var body struct {
		Type string `json:"type"`
	}

	var resp struct {
		Type     string  `json:"type"`
		Messages []int64 `json:"messages"`
	}

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	resp.Type = "read_ok"
	resp.Messages = data.getM()

	return n.Reply(msg, resp)
}

func topologyHandler(msg maelstrom.Message, n *maelstrom.Node) error {
	var body struct {
		Topology map[string][]string `json:"topology"`
		Type     string              `json:"type"`
	}
	var resp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	resp.Type = "topology_ok"
	// every node stores its neighbors
	// no need to store all data
	data.setNeighbors(body.Topology[n.ID()])

	return n.Reply(msg, resp)
}

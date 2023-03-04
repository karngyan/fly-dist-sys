package main

import (
	"encoding/json"
	"log"
	"sync"

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

func broadcastHandler(msg maelstrom.Message, n *maelstrom.Node) error {

	type req struct {
		Message int64  `json:"message"`
		MsgId   int64  `json:"msg_id"`
		Type    string `json:"type"`
	}

	var body req

	var resp struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// if message is not in the store, add it and broadcast it to neighbors
	if !data.getMByKey(body.Message) {
		data.addM(body.Message)
		neighbors := data.getNeighbours()
		for _, neighbor := range neighbors {
			if err := n.Send(neighbor, body); err != nil {
				log.Printf("Error sending message: %s", err)
			}
		}
	}

	if body.MsgId == 0 {
		// as we know this is a Send call from one of nodes
		return nil
	}

	resp.Type = "broadcast_ok"
	return n.Reply(msg, resp)
}

func broadcastRPCHandler(msg maelstrom.Message, n *maelstrom.Node, wg *sync.WaitGroup) error {
	wg.Done()
	return nil
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
	data.setNeighbours(body.Topology[n.ID()])

	return n.Reply(msg, resp)
}

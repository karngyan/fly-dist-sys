package main

import (
  "encoding/json"
  maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type HandlerFunc func(msg maelstrom.Message, n *maelstrom.Node) error

func handlerWithNode(h HandlerFunc, n *maelstrom.Node) maelstrom.HandlerFunc {
  return func(msg maelstrom.Message) error {
    return h(msg, n)
  }
}

var data = newStore()

func broadcastHandler(msg maelstrom.Message, n *maelstrom.Node) error {

  type req struct {
    MessageId int64  `json:"message"`
    Type      string `json:"type"`
    Hop       int    `json:"hop"`
  }

  var body req

  var resp struct {
    Type string `json:"type"`
  }

  if err := json.Unmarshal(msg.Body, &body); err != nil {
    return err
  }

  resp.Type = "broadcast_ok"
  data.addM(body.MessageId)

  nodeIds := n.NodeIDs()
  body.Hop = body.Hop + 1

  if body.Hop > 1 {
    return nil
  }

  for _, nodeId := range nodeIds {
    if nodeId == msg.Dest {
      continue
    }

    if err := n.Send(nodeId, body); err != nil {
      return err
    }
  }

  return n.Reply(msg, resp)
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
  for k, v := range body.Topology {
    data.addT(k, v)
  }

  return n.Reply(msg, resp)
}

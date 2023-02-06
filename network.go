package paxos

import (
	"log"
	"time"
)

type network struct {
	receiveQueue map[int]chan message
}

type nodeNetwork struct {
	id int
	network *network
}

func CreateNetwork(nodes ...int) *network {
	network := network{receiveQueue: make(map[int]chan message, 0)}

	for _, node := range nodes {
		network.receiveQueue[node] = make(chan message, 1024)
	}

	return &network
}

func (n *network) getNodeNetwork(id int) nodeNetwork {
	return nodeNetwork{id: id, network: n}
}

func (n *network) sendTo(message message) {
	log.Println("Send message from: ", message.from, " send to ", message.to, " value: ", message.value, " type: ", message.mtype)
	n.receiveQueue[message.to] <- message
}

func (n *network) receiveFrom(id int) *message {
	select {
	case retMessage := <- n.receiveQueue[id]:
		log.Println("Receive message from: ", retMessage.from, " send to ", retMessage.to, " value: ", retMessage.value, " type: ", retMessage.mtype)

		return &retMessage
	case <- time.After(time.Second):
		return nil
	}
}

func (nodeNetwork *nodeNetwork) send(message message) {
	nodeNetwork.network.sendTo(message)
} 

func (nodeNetwork *nodeNetwork) receive() *message {
	return nodeNetwork.network.receiveFrom(nodeNetwork.id)
}

package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
	"math"
)

// Global variables
const NODENUM := 2
var globalMap
var globalWG sync.WaitGroup
// Use a node structure to save information
type Node struct {
	id 	         int
	myNum        int
	highestNum   int
	replyTracker map[int]Node*   
	requestCS    bool
	defferedNode map[int]Node*
	nodeChannel  chan Message
}

type Message struct {
	messageType  string
	senderId     int
	receiver     Node*
	requestedNum int
}
// make a new node
func newNode(id int) Node* {
	n := Node{id, 0, 0, nil, false, map[int]*Node{}, make(chan Message)}
	return &n
}

// send a message
func (n Node*) sendMessage (msg Message) {
	receiver := msg.receiver
	fmt.Printf("[Node %d] Sending a <%s> message to Node %d at MemAddr %p \n", n.id,
		msg.messageType, receiver.id, globalMap[receiver.id])
	receiver.nodeChannel <- msg
}

// receive a message
func (n Node*) receiveMessage(msg Message) {
	n.highestNum = n.highestNum > msg.requestedNum ? n.highestNum: msg.requestedNum
	if !n.requestCS || (msg.requestedNum < n.myNum || (msg.requestedNum == n.myNum && msg.senderId < n.ID))
}

// main process
func (n Node*) mainProcess () {
	for {
		fmt.Println("Node %d enters non-critical section ", n.id)
		n.requestCS := true
		n.myNum := highestNum + 1
		// send messages to all other nodes
		for i := 0; i < NODENUM; i++ {
			if i == n.id{
				continue
			}
			msg := Message{"request", n.id, globalMap[i], n.myNum}
			go n.sendMessage(msg)
		}
	}
}

// receive process
func (n Node*) receiveProcess() {
}


// main process
func main () {
	globalMap := map[int]*Node{}
	globalWG.Add(NODENUM)
	fmt.Printf("Create Node now ...\n");
	for i := 0; i < NODENUM; i++ {
		new := newNode(i)
		globalMap[i] = new
	}
	fmt.Printf("Successfully created \d nodes.\n", NODENUM)

	for i := 0; i < NODENUM; i++ {
		go globalMap[i].mainProcess()
		go globalMap[i].receiveProcess()
		fmt.Println("Create main and receive process for node %d", i)
	}
	globalWG.wait()
	fmt.Println("All Nodes have entered entered and exited the Critical Section\nEnding programme now.\n")
} 

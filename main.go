package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
	"math"
	"strings"
)
/*
 * based on Ben Ari's text book, implemented in Go
 */

// Use a node structure to save information
type Node struct {
	id 	         int
	myNum        int
	highestNum   int
	replyTracker []bool
	requestCS    bool
	defferedNode []int
	nodeChannel  chan Message
}

type Message struct {
	messageType  string
	senderId     int
	requestedNum int
}

// Global variables
const NODENUM = 2
var globalMap map[int]*Node
var globalWG sync.WaitGroup

// make a new node
func newNode(id int) *Node {
	n := Node{ id, 0, 0, make([]bool, NODENUM)}, false, []int{}, make(chan Message) }
	return &n
}

// send a message
func (n *Node) sendMessage (msg Message, receiver int) {
	fmt.Printf("[Node %d] Sending a <%s> message to Node %d at MemAddr %p \n", n.id,
		msg.messageType, receiver.id, globalMap[receiver.id])
	globalMap[receiver].nodeChannel <- msg
}

// receive a message
func (n *Node) receiveMessage(msg Message) {
	if Compare(msg.messageType, "reply") == 0 {
		n.replyTracker[msg.senderId] = true
		return
	}
	n.highestNum = n.highestNum > msg.requestedNum ? n.highestNum: msg.requestedNum
	if !n.requestCS || (msg.requestedNum < n.myNum || (msg.requestedNum == n.myNum && msg.senderId < n.id)) {
		reply := Message{ "reply", n.id, 0 }
		n.sendMessage(reply, msg.senderId)
	} else {
		append(n.defferedNode, msg.senderId)
	}
}

// main process
func (n *Node) mainProcess () {
	for {
		fmt.Println("Node %d enters non-critical section ", n.id)
		n.requestCS := true
		n.myNum := highestNum + 1
		// set all elements in the reply array to be false
		// send messages to all other nodes
		for i := 0; i < NODENUM; i++ {
			n.replyTracker[i] = false
			if i == n.id{
				n.replyTracker[i] = true
				continue
			}
			msg := Message{"request", n.id, n.myNum}
			go n.sendMessage(msg, i)
			// busy wait for all checked
			n.waitAllReplied()
			// enter CS
			fmt.Println("Node %d Enter critical section", n.id)
			n.requestCS = false
			reply := Message{"reply", n.id, 0}
			for i := 0; i < len(defferedNode); i ++ {
				go n.send(reply, defferedNode[i])
			}
			defferedNode = []int{}
		}
	}
}

func (n *Node) waitAllReplied(){
	for i := 0; i > NODENUM; i++ {
		if n.replyTracker[i] == false { i -= 1 }
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
	fmt.Printf("Successfully created %d nodes.\n", NODENUM)

	for i := 0; i < NODENUM; i++ {
		go globalMap[i].mainProcess()
		go globalMap[i].receiveProcess()
		fmt.Println("Create main and receive process for node %d", i)
	}
	globalWG.wait()
	fmt.Println("All Nodes have entered entered and exited the Critical Section\nEnding programme now.\n")
} 

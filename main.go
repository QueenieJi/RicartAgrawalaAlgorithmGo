package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

/*
 * based on Ben Ari's text book, implemented in Go
 */

// Use a node structure to save information
type Node struct {
	id           int
	myNum        int
	highestNum   int
	replyTracker []bool
	requestCS    bool
	defferedNode []int
	nodeChannel  chan Message
	globalMap    map[int]*Node
}

type Message struct {
	messageType  string
	senderId     int
	requestedNum int
}

// Global variables
const NODENUM = 8

var globalWG sync.WaitGroup

// make a new node
func newNode(id int) *Node {
	n := Node{id, 0, 0, make([]bool, NODENUM), false, []int{}, make(chan Message), nil}
	return &n
}

// send a message
func (n *Node) sendMessage(msg Message, receiver int) {
	time.Sleep(1000 * time.Millisecond)
	fmt.Printf("[Node %d] Sending a <%s> message to Node %d at MemAddr %p \n", n.id,
		msg.messageType, receiver, n.globalMap[receiver])
	n.globalMap[receiver].nodeChannel <- msg
}

// receive a message
func (n *Node) receiveMessage() {
	msg := <-n.nodeChannel
	fmt.Printf("[Node %d] received a <%s> message from Node %d at MemAddr %p with request number %d\n", n.id,
		msg.messageType, msg.senderId, n.globalMap[msg.senderId], msg.requestedNum)
	if strings.Compare(msg.messageType, "reply") == 0 {
		n.replyTracker[msg.senderId] = true
		return
	}
	// fmt.Println("goodluck")
	if n.highestNum < msg.requestedNum {
		n.highestNum = msg.requestedNum
	}
	// if you put print here，就没有deadlock了
	if !n.requestCS || (msg.requestedNum < n.myNum || (msg.requestedNum == n.myNum && msg.senderId < n.id)) {
		reply := Message{"reply", n.id, 0}
		n.sendMessage(reply, msg.senderId)
	} else {
		n.defferedNode = append(n.defferedNode, msg.senderId)
		fmt.Printf("Node %d add node %d to its deferred node list\n", n.id, msg.senderId)
	}
}

// main process
func (n *Node) mainProcess() {
	for {
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Node %d enters non-critical section\n", n.id)
		n.requestCS = true
		n.myNum = n.highestNum + 1
		fmt.Printf("[node %d]mynum is %d\n", n.id, n.myNum)
		// set all elements in the reply array to be false
		// send messages to all other nodes
		for i := 0; i < NODENUM; i++ {
			n.replyTracker[i] = false
			if i == n.id {
				n.replyTracker[i] = true
				continue
			}
			msg := Message{"request", n.id, n.myNum}
			n.sendMessage(msg, i)
		}
		// busy wait for all checked
		n.waitAllReplied()
		// enter CS
		fmt.Printf("Node %d Enter critical section\n", n.id)
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Node %d Exit critical section\n", n.id)

		n.requestCS = false
		reply := Message{"reply", n.id, 0}
		for i := 0; i < len(n.defferedNode); i++ {
			n.sendMessage(reply, n.defferedNode[i])
		}
		n.defferedNode = []int{}
	}
}

func (n *Node) waitAllReplied() {
	for i := 0; i < NODENUM; i++ {
		if n.replyTracker[i] == false {
			i -= 1
		}
	}
}

// receive process
func (n *Node) receiveProcess() {
	for {
		n.receiveMessage()
	}
}

// main process
func main() {
	globalNodesMap := map[int]*Node{}
	globalWG.Add(NODENUM)
	fmt.Printf("Create Node now ...\n")
	for i := 0; i < NODENUM; i++ {
		globalNodesMap[i] = newNode(i)
		globalNodesMap[i].globalMap = globalNodesMap
	}
	fmt.Printf("Successfully created %d nodes.\n", NODENUM)

	for i := 0; i < NODENUM; i++ {
		go globalNodesMap[i].receiveProcess()
		fmt.Printf("Create receive process for node %d\n", i)
	}
	for i := 0; i < NODENUM; i++ {
		go globalNodesMap[i].mainProcess()
		fmt.Printf("Create main process for node %d\n", i)
	}
	globalWG.Wait()
	fmt.Println("All Nodes have entered entered and exited the Critical Section\nEnding programme now.\n")
}

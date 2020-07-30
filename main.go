package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Global variables
const NODENUM = 10

// Use a node structure to save information
type Node struct {
	myNum        int
	highestNum   int
	replyTracker    
	requestCS    bool
	deffered     [NODENUM]int

}

type Message struct {
	messageType  string
	senderId     int
	receiver     map[int]Node*
}
// make a new node
func newNode(id int) Node* {

}

// send a message
func (n Node*) sendMessage (msg Message) {

}

// receive a message
func (n Node*) receiveMessage(msg Message) {

}

// main process
func (n Node*) mainProcess () {

}

// receive process
func (n Node*) receiveProcess() {}


// main process
func main () {
	globalMap := map[int]*Node{}
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
	fmt.Println("All Nodes have entered entered and exited the Critical Section\nEnding programme now.\n")
} 

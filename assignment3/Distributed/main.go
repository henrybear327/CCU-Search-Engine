// netstat -vanp tcp | grep -e "50507"
// modified from https://confusedcoders.com/go/create-a-basic-distributed-system-in-go-lang-part-1

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// NodeInfo is the Information/Metadata about node
type NodeInfo struct {
	NodeID     int    `json:"nodeID"`
	NodeIPAddr string `json:"nodeIPAddr"`
	Port       string `json:"port"`
}

// Just for pretty printing the node info
func (node NodeInfo) String() string {
	return "NodeInfo: {\n\tnodeID:" + strconv.Itoa(node.NodeID) + ", \n\tnodeIPAddr:" + node.NodeIPAddr + ", \n\tport:" + node.Port + " \n}"
}

// AddToClusterMessage is a standard format for a Request/Response for adding node to cluster
type AddToClusterMessage struct {
	Source  NodeInfo `json:"source"`
	Dest    NodeInfo `json:"dest"`
	Message string   `json:"message"`
}

// Just for pretty printing Request/Response info
func (req AddToClusterMessage) String() string {
	return "AddToClusterMessage:{\nsource:" + req.Source.String() + ",\ndest: " + req.Dest.String() + ",\nmessage:" + req.Message + " }"
}

// praseCommandLine parses the provided parameters on command line
func praseCommandLine() (bool, string, string) {
	makeMasterOnError := flag.Bool("makeMasterOnError", false, "make this node master if unable to connect to the cluster ip provided.")
	clusterip := flag.String("clusterip", "127.0.0.1:8001", "ip address of any node to connnect")
	myport := flag.String("myport", "8001", "ip address to run this node on. default is 8001.")
	flag.Parse()

	return *makeMasterOnError, *clusterip, *myport
}

func check(info string, err error) {
	if err != nil {
		log.Fatalln("error on", info, ":", err)
	}
}

func setupNodes(clusterip, myport string) (int, NodeInfo, NodeInfo) {
	/* Generate id for myself */
	rand.Seed(time.Now().UTC().UnixNano())
	myID := rand.Intn(99999999)

	myIP, err := net.InterfaceAddrs()
	check("net.InterfaceAddrs()", err)
	me := NodeInfo{
		NodeID:     myID,
		NodeIPAddr: myIP[0].String(),
		Port:       myport,
	}

	destHost, destPort, err := net.SplitHostPort(clusterip)
	if err != nil {
		log.Fatal("net.SplitHostPort", err)
	}
	// log.Println("dest IP", destHost, "dest port", destPort)

	dest := NodeInfo{
		NodeID:     -1,
		NodeIPAddr: destHost,
		Port:       destPort,
	}
	log.Println("My details: ", me.String())

	return myID, me, dest
}

func getAddToClusterMessage(source, dest NodeInfo, message string) AddToClusterMessage {
	return AddToClusterMessage{
		Source: NodeInfo{
			NodeID:     source.NodeID,
			NodeIPAddr: source.NodeIPAddr,
			Port:       source.Port,
		},
		Dest: NodeInfo{
			NodeID:     dest.NodeID,
			NodeIPAddr: dest.NodeIPAddr,
			Port:       dest.Port,
		},
		Message: message,
	}
}

func connectToCluster(me, dest NodeInfo) bool {
	/* connect to this socket details provided */
	connOut, err := net.DialTimeout("tcp", net.JoinHostPort(dest.NodeIPAddr, dest.Port), time.Duration(10)*time.Second)
	if err != nil {
		log.Println("Couldn't connect to cluster.", me.NodeID)
		return false
	}

	log.Println("Connected to master. Sending a message to master.")
	text := "Please add node " + strconv.Itoa(me.NodeID) + " to master"
	requestMessage := getAddToClusterMessage(me, dest, text)
	json.NewEncoder(connOut).Encode(&requestMessage)

	decoder := json.NewDecoder(connOut)
	var responseMessage AddToClusterMessage
	decoder.Decode(&responseMessage)
	log.Println("Got response:\n" + responseMessage.String())

	return true
}

func handleConnection(connIn net.Conn, me NodeInfo) {
	var requestMessage AddToClusterMessage
	json.NewDecoder(connIn).Decode(&requestMessage)
	log.Println("Master got request:\n" + requestMessage.String())

	text := "Master received add slave request from " + strconv.Itoa(requestMessage.Source.NodeID) + " and added it"
	responseMessage := getAddToClusterMessage(me, requestMessage.Source, text)
	json.NewEncoder(connIn).Encode(&responseMessage)
	connIn.Close()
}

func listenOnPort(me NodeInfo) {
	/* Listen for incoming messages */
	ln, err := net.Listen("tcp", fmt.Sprint(":"+me.Port))
	check("net.Listen", err)
	defer ln.Close()
	/* accept connection on port */
	/* not sure if looping infinetely on ln.Accept() is good idea */
	for {
		connIn, err := ln.Accept()
		if err != nil {
			log.Println("Error received while listening.", me.NodeID)
			continue
		}
		log.Println("Accepted one connection")
		go handleConnection(connIn, me)
	}
}

func main() {
	makeMasterOnError, clusterip, myport := praseCommandLine()

	myID, me, dest := setupNodes(clusterip, myport)

	ableToConnect := connectToCluster(me, dest)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		sig := <-sigs
		log.Println(sig)
		os.Exit(0)
	}()

	if ableToConnect || (!ableToConnect && makeMasterOnError) {
		if makeMasterOnError {
			log.Println("Will start this node as master.")
			me.NodeID = 0 // master ID = 0
		}

		listenOnPort(me)
	} else {
		log.Println("Quitting system. Set makeMasterOnError flag to make the node master.", myID)
	}
}

package main

import (
	"bufio"
	"log"

	"gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
)

var (
	rcv_from_client = make(chan string)
	send_to_client  []chan string
	client_list     []Client
)

// Function to handle incoming connection from client
func handleClientStream(stream net.Stream) {

	buff := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	send_chan := make(chan string)
	send_to_client = append(send_to_client, send_chan)

	log.Println("Incoming connection from client", stream.Conn().RemotePeer().Pretty())

	go receiveDataFromClient(buff, send_chan)
	go sendDataToClient(buff, send_chan)
}

// Goroutine that send data received from the stream to the command handler
func receiveDataFromClient(buff *bufio.ReadWriter, response_chan chan string) {

	for {
		data, _ := buff.ReadString('\n')
		unmarshalCommand(data, response_chan)
	}
}

// Goroutine that write data received from a channel to the stream
func sendDataToClient(buff *bufio.ReadWriter, send_chan chan string) {

	for {
		data := <-send_chan

		buff.WriteString(data)
		buff.Flush()
	}
}

// Send data to clients except ones in exceptions
func broadCastToClients(data string, exceptions ...chan string) {

	for _, client := range send_to_client {
		skip := false
		for _, exception := range exceptions {
			if exception == client {
				skip = true
			}
		}
		if !skip {
			client <- data
		}
	}
}

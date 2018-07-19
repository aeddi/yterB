package main

import (
	"bufio"

	"gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
)

var (
	send_to_relay []chan string
	client        Client
	client_list   []Client
)

// Function to handle incoming connection from relay
func handleRelayStream(stream net.Stream) {

	buff := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	send_chan := make(chan string)
	send_to_relay = append(send_to_relay, send_chan)

	consoleLog("Incoming connection from relay " + stream.Conn().RemotePeer().Pretty())

	go receiveDataFromRelay(buff, send_chan)
	go sendDataToRelay(buff, send_chan)
}

// Goroutine that send data received from the stream to the command handler
func receiveDataFromRelay(buff *bufio.ReadWriter, response_chan chan string) {

	for {
		data, _ := buff.ReadString('\n')
		unmarshalCommand(data, response_chan)
	}
}

// Goroutine that write data received from a channel to the stream
func sendDataToRelay(buff *bufio.ReadWriter, send_chan chan string) {

	for {
		data := <-send_chan

		buff.WriteString(data)
		buff.Flush()
	}
}

// Send data to relays except ones in exceptions
func broadCastToRelays(data string, exceptions ...chan string) {

	for _, relay := range send_to_relay {
		skip := false
		for _, exception := range exceptions {
			if exception == relay {
				skip = true
			}
		}
		if !skip {
			relay <- data
		}
	}
}

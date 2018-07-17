package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
)

var rcv_from_relay = make(chan string)
var send_to_relay []chan string

// Function to handle incoming connection
func handleStream(stream net.Stream) {

	buff := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	send_chan := make(chan string)
	send_to_relay = append(send_to_relay, send_chan)

	log.Println("Incoming connection from peer", stream.Conn().RemotePeer().Pretty())

	go sendData(buff, send_chan)
	go receiveData(buff)
}

// Goroutine that write data received from the stream to a channel
func receiveData(buff *bufio.ReadWriter) {

	for {
		str, _ := buff.ReadString('\n')
		rcv_from_relay <- str
	}
}

// Goroutine that write data received from a channel to the stream
func sendData(buff *bufio.ReadWriter, send_chan chan string) {

	for {
		str := <-send_chan

		buff.WriteString(str)
		buff.Flush()
	}
}

// Goroutine: tmp function to test other the program (print received data)
func readMessage() {

	for {
		str := <-rcv_from_relay

		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf(str)
		}
	}
}

// Goroutine: tmp function to test other the program (send inputed data)
func writeMessage() {

	std_reader := bufio.NewReader(os.Stdin)

	for {
		send_data, err := std_reader.ReadString('\n')

		if err != nil {
			removeAddressFile(1)
		}

		for _, send_chan := range send_to_relay {
			send_chan <- fmt.Sprintf("%s\n", send_data)
		}
	}
}

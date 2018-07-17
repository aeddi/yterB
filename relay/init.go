package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"strconv"
)

func main() {

	port := flag.Int("p", 0, "Specify an unused port")
	flag.Parse()

	port_str := strconv.Itoa(*port)
	checkIfPortIsAvailable(port_str)
	peer_addresses := getAddressesFromFiles()
	relay_host := createRelayHost(port_str)

	// host.SetStreamHandler("/ytreb/relay_to_relay", handleRelayStream)
	// host.SetStreamHandler("/ytreb/client_to_relay", handleClientStream)
	relay_host.SetStreamHandler("/ytreb/relay2relay", handleStream)

	for _, peer_address := range peer_addresses {
		peer_id := addAddrToPeerstore(relay_host, peer_address)
		stream, err := relay_host.NewStream(context.Background(), peer_id, "/ytreb/relay2relay")

		log.Println("Init connection with peer", peer_id.Pretty())

		if err == nil {
			buff := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			send_chan := make(chan string)
			send_to_relay = append(send_to_relay, send_chan)

			log.Println("Connected to peer", peer_id.Pretty())
			go sendData(buff, send_chan)
			go receiveData(buff)
		} else {
			log.Println("Can't connect to peer", peer_id.Pretty())
		}
	}

	go writeMessage()
	go readMessage()
	select {}
}

package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"strconv"
)

func main() {

	port := flag.Int("p", 0, "Specify an unused port greater than 1023")
	flag.Parse()

	checkIfPortIsAvailable(*port)
	peer_addresses := getAddressesFromFiles()
	relay_host := createRelayHost(strconv.Itoa(*port))

	relay_host.SetStreamHandler("/ytreb/relay2relay", handleRelayStream)
	relay_host.SetStreamHandler("/ytreb/client2relay", handleClientStream)

	for _, peer_address := range peer_addresses {
		peer_id := addAddrToPeerstore(relay_host, peer_address)
		stream, err := relay_host.NewStream(context.Background(), peer_id, "/ytreb/relay2relay")

		log.Println("Init connection with peer", peer_id.Pretty())

		if err == nil {
			buff := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			send_chan := make(chan string)
			send_to_relay = append(send_to_relay, send_chan)

			log.Println("Connected to peer", peer_id.Pretty())
			go receiveDataFromRelay(buff, send_chan)
			go sendDataToRelay(buff, send_chan)
		} else {
			log.Println("Can't connect to peer", peer_id.Pretty())
		}
	}

	select {}
}

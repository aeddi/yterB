package main

import (
	"bufio"
	"context"
	"crypto/x509"
	"log"
)

func initClient(name string) {

	key_pair := generateKeyPair()
	relay_addresses := getAddressesFromFiles()

	if len(relay_addresses) == 0 {
		log.Fatalln("Error: no relay available")
	}

	public_key, _ := x509.MarshalPKIXPublicKey(&key_pair.PublicKey)
	client.Name = name
	client.Public_key = string(public_key)
	client_host := createRelayHost(getAvailablePort())

	client_host.SetStreamHandler("/ytreb/client2relay", handleRelayStream)

	for _, relay_address := range relay_addresses {
		peer_id := addAddrToPeerstore(client_host, relay_address)
		stream, err := client_host.NewStream(context.Background(), peer_id, "/ytreb/client2relay")

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

	// Send to all relays the client identity (public key, name, id and address)
	broadCastToRelays(marshalCommand("client_register", client))
	select {}
}

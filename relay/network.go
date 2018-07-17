package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	netstd "net"

	"github.com/libp2p/go-libp2p"

	"gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	"gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

// Generate RSA-2048 key pair
func generateKeyPair() crypto.PrivKey {

	private_key, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

	if err != nil {
		log.Fatal("Error during key pair generation: ", err)
	}

	log.Println("Key pair generetad successfully")

	return private_key
}

// Check if TCP port is available before using it
func checkIfPortIsAvailable(port string) {

	listen, err := netstd.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalln("Error: you must specify an unused port with -p flag\n", err)
	}

	log.Println("Port", port, "is available")
	listen.Close()
}

// Create a libp2p host
func createRelayHost(port string) host.Host {

	private_key := generateKeyPair()
	multi_addr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/" + port)
	relay_host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(multi_addr),
		libp2p.Identity(private_key),
	)

	if err != nil {
		log.Fatalln("Error during host creation", err)
	}

	createAddressFile(port, "/ip4/127.0.0.1/tcp/"+port+"/ipfs/"+relay_host.ID().Pretty())
	log.Println("Relay host successfully created with id:", relay_host.ID().Pretty())
	return relay_host
}

// Add a remote address to the peerstore to be able to communicate with the peer
func addAddrToPeerstore(relay_host host.Host, peer_address string) peer.ID {

	multi_addr, err := multiaddr.NewMultiaddr(peer_address)
	if err != nil {
		log.Println("Error during peer multiaddress creation:", err)
		removeAddressFile(1)
	}
	peer_infos, err := multi_addr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		log.Println("Error during peer infos obtention:", err)
		removeAddressFile(1)
	}
	peer_id, err := peer.IDB58Decode(peer_infos)
	if err != nil {
		log.Println("Error during peer ID decoding:", err)
		removeAddressFile(1)
	}

	target_addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peer_id)))
	target_addr_dec := multi_addr.Decapsulate(target_addr)

	relay_host.Peerstore().AddAddr(peer_id, target_addr_dec, peerstore.PermanentAddrTTL)
	log.Println("Relay peer ID added to peerstore")
	return peer_id
}

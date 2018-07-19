package main

import (
	"context"
	"crypto/rand"
	"fmt"
	netstd "net"
	"os"
	"strconv"

	"github.com/libp2p/go-libp2p"

	"gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	"gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

// Generate RSA-2048 key pair for P2P communication
func generateKeyPairP2P() crypto.PrivKey {

	private_key, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

	if err != nil {
		consoleLog("Error during key pair generation: " + err.Error())
		os.Exit(1)
	}

	consoleLog("Key pair generetad successfully")

	return private_key
}

// Check if TCP port is available before using it
func checkIfPortIsAvailable(port int) bool {

	listen, err := netstd.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		return false
	}

	listen.Close()
	return true
}

// Search an available TCP port by iterating from 1024 to 65535
func getAvailablePort() string {

	for port := 1024; port <= 65535; port++ {
		if checkIfPortIsAvailable(port) {
			return strconv.Itoa(port)
		}
	}

	consoleLog("Error: unable to find an available TCP port")
	os.Exit(1)
	return ""
}

// Create a libp2p host
func createRelayHost(port string) host.Host {

	private_key_p2p := generateKeyPairP2P()
	multi_addr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/" + port)
	relay_host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(multi_addr),
		libp2p.Identity(private_key_p2p),
	)
	client.Peer_id = relay_host.ID().Pretty()
	client.Address = multi_addr.String()

	if err != nil {
		consoleLog("Error during host creation " + err.Error())
		os.Exit(1)
	}

	consoleLog("Relay host successfully created with id: " + client.Peer_id)
	return relay_host
}

// Add a remote address to the peerstore to be able to communicate with the peer
func addAddrToPeerstore(relay_host host.Host, peer_address string) peer.ID {

	multi_addr, err := multiaddr.NewMultiaddr(peer_address)
	if err != nil {
		consoleLog("Error during peer multiaddress creation: " + err.Error())
		os.Exit(1)
	}
	peer_infos, err := multi_addr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		consoleLog("Error during peer infos obtention: " + err.Error())
		os.Exit(1)
	}
	peer_id, err := peer.IDB58Decode(peer_infos)
	if err != nil {
		consoleLog("Error during peer ID decoding: " + err.Error())
		os.Exit(1)
	}

	target_addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peer_id)))
	target_addr_dec := multi_addr.Decapsulate(target_addr)

	relay_host.Peerstore().AddAddr(peer_id, target_addr_dec, peerstore.PermanentAddrTTL)
	consoleLog("Relay peer ID added to peerstore")
	return peer_id
}

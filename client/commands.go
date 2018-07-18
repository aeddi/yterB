package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
)

type Command struct {
	Header  string
	Command string
}

type Client struct {
	Name       string
	Public_key string
	Peer_id    string
	Address    string
}

type Authcode struct {
	Authcode  string
	MessageID string
	Sender    string
}

type Message struct {
	Content   string
	Datetime  string
	Authcode  string
	Sender    string
	Recipient string
	MessageID string
}

var message_queue []Message

// Create a Command struct then encode it to a json string
func marshalCommand(header string, data interface{}) string {

	var command Command

	encoded, _ := json.Marshal(data)
	command.Header = header
	command.Command = string(encoded)
	encoded, _ = json.Marshal(command)

	return string(encoded) + "\n"
}

// Decode json data and call the function associated to header
func unmarshalCommand(data string, response_chan chan string) {

	var command Command

	json.Unmarshal([]byte(data), &command)
	switch command.Header {
	case "client_register":
		var client Client
		json.Unmarshal([]byte(command.Command), &client)
		registerClient(client)
	}
	////////////////////////////
	_ = response_chan
	////////////////////////////
}

// Register client in databse then add contact to the GUI
func registerClient(client Client) {

	registered := false

	log.Println("Command client_register received")
	for _, client_registered := range client_list {
		if client_registered.Peer_id == client.Peer_id {
			registered = true
			break
		}
	}

	if registered {
		log.Println("Client", client.Name, "registered already")
		return
	} else {
		client_list = append(client_list, client)
		log.Println("Registering client with name:", client.Name)
		addContactToGUI(client)
	}
}

// Init message sending by requesting an auth code
func requestAuthcode(message Message) {

	var request Authcode
	hasher := sha1.New()
	hasher.Write([]byte(message.Content + message.Datetime + message.Recipient))
	message.MessageID = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	request.MessageID = message.MessageID
	request.Sender = client.Peer_id

	message_queue = append(message_queue, message)

	send_to_relay[rand.Intn(len(send_to_relay))] <- marshalCommand("authcode_request", request)
}

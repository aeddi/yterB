package main

import (
	"encoding/json"
	"log"
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
		registerClient(client, response_chan)
	case "authcode_request":
		var request Authcode
		json.Unmarshal([]byte(command.Command), &request)
		authcodeResponse(request)
	}
}

func authcodeResponse(request Authcode) {

	log.Println("Authcode request received from", request.Sender, "for message", request.MessageID)
}

func registerClient(client Client, response_chan chan string) {

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
		log.Println("Registering client with name:", client.Name)
		log.Println("Sending client identity to other clients")
		broadCastToClients(marshalCommand("client_register", client), response_chan)

		log.Println("Sending other client identities to the new client")
		for _, client_registered := range client_list {
			response_chan <- marshalCommand("client_register", client_registered)
		}
		client_list = append(client_list, client)
	}
}

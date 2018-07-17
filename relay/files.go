package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
)

const tmp_folder = "/tmp/ytreb"

var address_file string

// Catch interrupt and kill signals to delete address file before termination
func cleanOnKill() {

	signal_chan := make(chan os.Signal, 1)

	signal.Notify(signal_chan, os.Interrupt, os.Kill)
	go func() {
		signal := <-signal_chan
		log.Println("Signal", signal, "received! Cleaning tmp file...")
		removeAddressFile(0)
	}()
}

// Remove address file from tmp folder
func removeAddressFile(error_code int) {

	err := os.Remove(address_file)

	if err != nil {
		log.Fatalln("Error during tmp file deletion:", err)
	}

	os.Exit(error_code)
}

// Create tmp folder to store adress files
func createTmpFolder() {

	err := os.MkdirAll(tmp_folder, 0700)

	if err != nil {
		log.Fatalln("Error during tmp folder creation:", err)
	}

	log.Println("Tmp folder created:", tmp_folder)
}

// Create a file in tmp folder containing an address
func createAddressFile(filename string, address string) {

	address_file = path.Join(tmp_folder, filename)
	err := ioutil.WriteFile(address_file, []byte(address), 0600)

	if err != nil {
		log.Fatalln("Error during address file creation:", err)
	}

	log.Println("Address file created:", address_file)
	cleanOnKill()
}

// Return an array of addresses got from all files contained in tmp folder
func getAddressesFromFiles() []string {

	var addresses []string

	createTmpFolder()

	err := filepath.Walk(tmp_folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			address, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			addresses = append(addresses, string(address))
		}
		return nil
	})

	if err != nil {
		log.Fatalln("Error during port files listing:", err)
	}

	log.Printf("Found %d other relay(s)\n", len(addresses))
	return addresses
}

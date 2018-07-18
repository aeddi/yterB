package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const tmp_folder = "/tmp/ytreb"

// Create tmp folder to store adress files
func createTmpFolder() {

	err := os.MkdirAll(tmp_folder, 0700)

	if err != nil {
		log.Fatalln("Error during tmp folder creation:", err)
	}

	log.Println("Tmp folder created:", tmp_folder)
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

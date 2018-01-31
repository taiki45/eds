package main

import (
	"flag"
	"log"

	//"github.com/taiki45/eds/client"
	"github.com/taiki45/eds/server"
)

func main() {
	flag.Parse()
	command := flag.Arg(0)

	if command == "server" {
		log.Printf("Running server...")
		server.Run()
	} else if command == "client" {
		log.Printf("Running client...")
		//client.Run()
	} else {
		flag.PrintDefaults()
		log.Fatal("Invalid command line arguments.")
	}
}

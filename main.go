package main

import "fmt"

import "time"
import "github.com/hashicorp/serf/client"
import "p2p"

func main() {
	client, _ := client.NewRPCClient("127.0.0.1:7373")
	members, _ := client.Members()
	for {
		for member := range members {
			fmt.Println("Found member", member)
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Nodename: ", FOURTYWO)
}

package kyoo

import "fmt"

import "time"
import "github.com/hashicorp/memberlist"
import "flag"
import "os"

const (
	GetNodeName uint64 = iota
)

func Kyoo() int {
	fmt.Println("Starting kyoo")
	var bindIp = flag.String("bindip", "127.0.0.1", "IP To bind to")
	hostname, _ := os.Hostname()
	var nodeName = flag.String("nodename", hostname, "NodeName")

	flag.Parse()
	config := memberlist.DefaultLocalConfig()
	config.BindAddr = *bindIp
	config.Name = *nodeName
	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}
	list.Join([]string{"127.0.0.1"})

	/*
		serfClient, _ := client.NewRPCClient("127.0.0.1:7373")
		// client.Respond(GetNodeName, []byte{1, 2, 3})
		// params := client.QueryParam{}
		// serfClient.query(queryParam)
		ci := make(chan map[string]interface{})
		serfClient.Stream("", ci)
		go func() {
			for {
				foo := <-ci
				fmt.Println("Received: ", foo)
			}
		}()
		for {
			members, _ := serfClient.Members()
			for member := range members {
				fmt.Println("Found member", member)
			}
			time.Sleep(1 * time.Second)
		}
	*/
	for {

		for _, member := range list.Members() {
			fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
		}
		time.Sleep(1 * time.Second)
	}
	return 0

}

package kyoo

import "fmt"

import "flag"
import "os"
import "kyoo/memberview"
import "os/signal"

var stopChan = make(chan int)

const (
	GetNodeName uint64 = iota
)

func Kyoo() int {
	fmt.Println("Starting kyoo")

	hostname, _ := os.Hostname()
	var nodeName = flag.String("nodename", hostname, "NodeName")

	flag.Parse()

	memberView := memberview.Create(*nodeName)
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
	//	eventStr := strings.Join([]string{"UserEvent", *nodeName}, " ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	memberview.Shutdown(memberView)
	return 0

}

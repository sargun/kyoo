package memberview

import "github.com/hashicorp/serf/serf"
import "time"
import "bytes"
import "math/rand"
import "strings"
import "log"
import "fmt"

type member struct {
}
type NodeStateValue interface {
}
type question struct {
}

//type NodeState {
//	attribs map[int]NodeStateValue
//}

type MemberView struct {
	eventChannel chan serf.Event
	qChan        chan question
	serf         *serf.Serf
	logger       *log.Logger
	members      map[string]member
	// TODO: Generalize this into a command channel
	shutDownChannel chan (chan int)
}

func handleMemberEvent(memberView *MemberView, memberEvent serf.MemberEvent) {

	for member := range memberEvent.Members {
		switch memberEvent.Type {
		case serf.EventMemberJoin:
			fmt.Println("Member join event:", memberEvent.Members[member].Name, memberEvent.String())
		default:
			fmt.Println("Received unknown member event: %v %v", memberEvent.Members[member].Name, memberEvent.String())
		}
	}
}

func handleUserEvent(memberView *MemberView, userEvent serf.UserEvent) {

}
func handleEvent(memberView *MemberView, event serf.Event) {
	switch event.(type) {
	case serf.UserEvent:
		ue, ok := event.(serf.UserEvent)
		if !ok {
			memberView.logger.Panic("Unable to convert to user event.")
		}
		handleUserEvent(memberView, ue)
	case serf.MemberEvent:
		me, ok := event.(serf.MemberEvent)
		if !ok {
			memberView.logger.Panic("Unable to convert to member event.")
		}
		handleMemberEvent(memberView, me)
	}
	//	fmt.Println("Got serf event, ", *event)
}
func handleTick(memberView *MemberView) {
	if len(memberView.serf.Members()) == 1 {
		memberView.serf.Join([]string{"127.0.0.1:7946"}, false)
	}
}
func loop(memberView *MemberView) {
	defer func() {
		memberView.serf.Leave()
		memberView.serf.Shutdown()
	}()
	tick := time.Tick(1 * time.Second) // TODO: Make configurable
	memberView.logger.Println("Joining Serf Cluster")
	for {
		select {
		case event := <-memberView.eventChannel:
			handleEvent(memberView, event)
		case <-tick:
			handleTick(memberView)
		case syncChan := <-memberView.shutDownChannel:
			memberView.serf.Leave()
			memberView.serf.Shutdown()
			syncChan <- 0
			return
		}
	}
}

func Create(nodeName string) *MemberView {
	memberView := new(MemberView)
	memberView.shutDownChannel = make(chan (chan int), 1)
	var buf bytes.Buffer
	memberView.logger = log.New(&buf, "Memberview: ", log.Lshortfile)
	snapshotPath := strings.Join([]string{"/tmp/serf-", nodeName}, "")

	config := serf.DefaultConfig()
	config.NodeName = nodeName
	config.SnapshotPath = snapshotPath
	config.MemberlistConfig.Name = nodeName

	port := rand.Intn(16384) + 16384
	config.MemberlistConfig.BindPort = port
	config.MemberlistConfig.AdvertisePort = port

	memberView.eventChannel = make(chan serf.Event, 1)
	memberView.qChan = make(chan question, 1)
	memberView.members = make(map[string]member)
	config.EventCh = memberView.eventChannel
	var err error
	memberView.serf, err = serf.Create(config)
	if err != nil {
		panic(err)
	}
	go loop(memberView)
	return memberView
}
func Shutdown(memberView *MemberView) {
	syncChan := make(chan int)
	memberView.shutDownChannel <- syncChan
	<-syncChan
}

/*
	for {
		fmt.Println("---------------")
		members := s.Members()
		//s.UserEvent(eventStr, []byte{0, 1, 2, 3}, false)
		for n := range members {
			member := members[n]
			fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
		}
		time.Sleep(1 * time.Second)
	}
*/
/*	go func() {
	for {
		Event := <-eventChannel
		if Event.EventType() == serf.EventMemberJoin {
			fmt.Println("HI: ", Event.String())
		}
		fmt.Println("Received channel: ", Event)
	}
}()*/

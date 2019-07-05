package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/gofrs/uuid"
	stan "github.com/nats-io/stan.go"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		url        = flag.String("url", stan.DefaultNatsURL, "NATS Server URLs, separated by commas")
		clusterID  = flag.String("cluster_id", "store3", "Cluster ID")
		clientID   = flag.String("client_id", "", "Client ID")
		queueGroup = flag.String("queue-group", "", "Queue group ID")
	)
	flag.Parse()

	if *clientID == "" {
		*clientID = uuid.Must(uuid.NewV4()).String()
	}
	sc, err := stan.Connect(*clusterID, *clientID,
		stan.NatsURL(*url),
		stan.Pings(10, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("Connection lost: %v", reason)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	sub, err := sc.QueueSubscribe("ECHO", *queueGroup, func(msg *stan.Msg) {
		log.Printf("%10s | %s\n", msg.Subject, string(msg.Data))
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Fatal(err)
	}

	// Wait for Ctrl+C
	doneCh := make(chan bool)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh
		sub.Unsubscribe()
		doneCh <- true
	}()
	<-doneCh
}

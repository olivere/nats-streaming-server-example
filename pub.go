package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofrs/uuid"
	stan "github.com/nats-io/stan.go"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		url       = flag.String("url", stan.DefaultNatsURL, "NATS Server URLs, separated by commas")
		clusterID = flag.String("cluster_id", "store3", "Cluster ID")
		clientID  = flag.String("client_id", "", "Client ID")
	)
	flag.Parse()

	if *clientID == "" {
		*clientID = uuid.Must(uuid.NewV4()).String()
	}

	// Connect to NATS Streaming Server cluster
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

	// Publish some messages, synchronously
	var i int64
	for {
		i++
		now := time.Now().Format(time.RFC3339)
		payload := fmt.Sprintf("%08d %s", i, now)
		err := sc.Publish("ECHO", []byte(payload))
		if err != nil {
			log.Fatal(err)
		}
		// Sleep for a random time of up to 1s
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	}
}

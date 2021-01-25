package main

import (
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	// "go.etcd.io/etcd/client/v3/experimental/recipes"

	// "fmt"
	"log"
	"time"
)

func main() {
	log.Println("CREATE CLIENT")
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"etcd-0:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	log.Println("CREATE SESSION")
	//send TTL updates to server each 1s. If failed to send (client is down or without communications), lock will be released
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(1))
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	start := time.Now()
	l1 := NewRWMutex(session, "lock1")
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("time to create lock object:", elapsed)

	start = time.Now()
	l1.RLock()
	t = time.Now()
	elapsed = t.Sub(start)
	log.Println("time to acquire lock :", elapsed)
	log.Println("READ LOCK AQCUIRED")

	time.Sleep(5 * time.Second)

	start = time.Now()
	l1.RUnlock()
	t = time.Now()
	elapsed = t.Sub(start)
	log.Println("time to release lock :", elapsed)
	log.Println("READ LOCK RELEASED")

	l2 := NewRWMutex(session, "lock1")

	start = time.Now()
	l2.Lock()
	t = time.Now()
	elapsed = t.Sub(start)
	log.Println("time to acquire lock:", elapsed)
	log.Println("WRITE LOCK AQCUIRED")

	time.Sleep(5 * time.Second)

	start = time.Now()
	l2.Unlock()
	t = time.Now()
	elapsed = t.Sub(start)
	log.Println("time to release lock:", elapsed)
	log.Println("WRITE LOCK RELEASED")
}

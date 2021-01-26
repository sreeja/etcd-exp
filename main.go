package main

import (
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	// "go.etcd.io/etcd/client/v3/experimental/recipes"
	"github.com/sreeja/etcd-exp/rwlock"

	// "fmt"
	"log"
	"os"
	"time"
)

func main() {
	filename := "log" + os.Args[1]
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println("CREATE CLIENT")
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"etcd-" + os.Args[1] + ":2379"}})
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
	l1 := rwlock.NewRWMutex(session, "lock1")
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("time to create lock object:", elapsed)

	for i := 0; i <= 1000; i++ {
		start = time.Now()
		rlerr := l1.RLock()
		if rlerr != nil {
			log.Fatal(rlerr)
		}
		t = time.Now()
		elapsed = t.Sub(start)
		log.Println("time to acquire read lock :", elapsed)
		// log.Println("READ LOCK AQCUIRED")

		time.Sleep(5 * time.Millisecond)

		start = time.Now()
		urlerr := l1.RUnlock()
		if urlerr != nil {
			log.Fatal(urlerr)
		}
		t = time.Now()
		elapsed = t.Sub(start)
		log.Println("time to release read lock :", elapsed)
		// log.Println("READ LOCK RELEASED")
		time.Sleep(5 * time.Millisecond)
	}

	l2 := rwlock.NewRWMutex(session, "lock1")

	for i := 0; i <= 1000; i++ {
		start = time.Now()
		lerr := l2.Lock()
		if lerr != nil {
			log.Fatal(lerr)
		}
		t = time.Now()
		elapsed = t.Sub(start)
		log.Println("time to acquire lock:", elapsed)
		// log.Println("WRITE LOCK AQCUIRED")

		time.Sleep(5 * time.Millisecond)

		start = time.Now()
		ulerr := l2.Unlock()
		if ulerr != nil {
			log.Fatal(ulerr)
		}
		t = time.Now()
		elapsed = t.Sub(start)
		log.Println("time to release lock:", elapsed)
		// log.Println("WRITE LOCK RELEASED")
		time.Sleep(5 * time.Millisecond)
	}
}

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/cretz/bine/tor"
)

func main() {
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	t, err := tor.Start(nil, nil)
	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}
	defer t.Close()

	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	// Create an onion service to listen on any port but show as 80
	onion, err := t.Listen(listenCtx, &tor.ListenConf{RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Unable to create onion service: %v", err)
	}
	defer onion.Close()
	fmt.Printf("%v.onion\n", onion.ID)

	if _, err := os.Stat("./tor.url"); err != nil {
		if os.IsNotExist(err) {
			os.Create("./tor.url")
			ioutil.WriteFile("./tor.url", []byte(onion.ID), 755)
		}
	} else {
		ioutil.WriteFile("./tor.url", []byte(onion.ID), 755)
	}

	//fmt.Println("Press enter to exit")
	// Serve the current folder from HTTP
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		http.Serve(onion, http.FileServer(http.Dir("./upload/")))
	}()

	wg.Wait()
}

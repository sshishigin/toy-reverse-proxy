package main

import (
	"fmt"
	"log"
	"net/http"
	"toy-reverse-proxy/src/buckets"
)

func main() {
	serverBucket := buckets.GetServerBucket("simple")
	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		serverBucket.Do(rw, req)
		return
	})
	fmt.Println("[Listening and serving on localhost:8080]")
	log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}

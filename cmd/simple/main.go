package main

import (
	"log"
	"net/http"
	"toy-reverse-proxy/src/simple"
)

func main() {
	serverBucket := simple.NewSimpleServerBucket()
	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		serverBucket.Do(rw, req)
		return
	})
	log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}

package simple

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type ServerBucket struct {
	ServerList       []simpleServer
	pointer          int
	serversAvailable []*simpleServer
	maxRetries       int
}

func (sb *ServerBucket) getServer() *simpleServer {
	server := sb.serversAvailable[sb.pointer]
	if sb.pointer < len(sb.serversAvailable)-1 {
		sb.pointer++
	} else {
		sb.pointer = 0
	}
	if server.available {
		return server
	}

	return nil
}

func (sb *ServerBucket) Do(rw http.ResponseWriter, req *http.Request) {
	for i := 0; i < sb.maxRetries; i++ {
		server := sb.getServer()
		if server == nil {
			continue
		}
		req.Host = server.Location.Host
		req.URL.Host = server.Location.Host
		req.URL.Scheme = server.Location.Scheme
		req.RequestURI = ""
		response, err := http.DefaultClient.Do(req)

		if err != nil {
			_, _ = fmt.Fprintln(rw, err)
			log.Print("request failed\n")
			return
		}
		if response.StatusCode == 200 || (response.StatusCode != 429 && response.StatusCode < 500) {
			server.fails = 0
			rw.WriteHeader(response.StatusCode)
			_, err = io.Copy(rw, response.Body)
			if err != nil {
				log.Println(err)
			}
			log.Printf("successful response from %s with status %d\n", req.URL.Host, response.StatusCode)
			return
		}

		if server.fails == server.maxFails {
			go server.excludeWithTimeout()
		} else {
			server.fails++
		}
	}
	log.Printf("failed by all %d servers\n", len(sb.serversAvailable))
	rw.WriteHeader(500)
}

func NewSimpleServerBucket() (sb ServerBucket) {
	var servers []simpleServer
	var serversAvailable []*simpleServer
	serverFile, _ := os.OpenFile("server-list", os.O_RDONLY, 0666)
	reader := bufio.NewScanner(serverFile)
	for reader.Scan() {
		directives := strings.Split(reader.Text(), ` `)
		address := directives[0]
		timeout, err := strconv.Atoi(directives[2])
		if err != nil {
			log.Fatal(err)
		}
		maxFails, err := strconv.Atoi(directives[3])
		if err != nil {
			log.Fatal(err)
		}
		host, err := url.Parse(address)
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, simpleServer{host, true, time.Duration(timeout) * time.Millisecond, maxFails, 0})
	}
	for i := range servers {
		serversAvailable = append(serversAvailable, &servers[i])
	}
	fmt.Printf("%d servers found in config \n", len(serversAvailable))
	sb = ServerBucket{servers, 0, serversAvailable, 3}
	return
}

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
	for _ = range sb.serversAvailable {
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
			fmt.Printf("%s request failed\n", time.Now().Format("2006-01-02 15:04:05"))
			return
		}
		if response.StatusCode == 200 || (response.StatusCode != 429 && response.StatusCode < 500) {
			rw.WriteHeader(response.StatusCode)
			io.Copy(rw, response.Body)
			fmt.Printf("[%s] successful response from %s with status %d\n", time.Now().Format("2006-01-02 15:04:05"), req.URL.Host, response.StatusCode)
			return
		}
		server.excludeWithTimeout()
	}
	fmt.Printf("%s failed by all %d servers\n", time.Now().Format("2006-01-02 15:04:05"), len(sb.serversAvailable))
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
		host, err := url.Parse(address)
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, simpleServer{host, true, time.Duration(timeout) * time.Second})
	}
	for i := range servers {
		serversAvailable = append(serversAvailable, &servers[i])
	}
	fmt.Printf("%d servers found in config \n", len(serversAvailable))
	sb = ServerBucket{servers, 0, serversAvailable}
	return
}

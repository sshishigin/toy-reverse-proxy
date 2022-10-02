package weighted

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

type WeightedServerBucket struct {
	ServerList       []weightedServer
	pointer          int
	serversAvailable []*weightedServer
}

func (sb *WeightedServerBucket) getServer() *weightedServer {
	server := sb.serversAvailable[sb.pointer]
	if sb.pointer < len(sb.serversAvailable)-1 {
		sb.pointer++
	} else {
		sb.pointer = 0
	}
	for {
		if server.available {
			return server
		}
	}

}
func (sb *WeightedServerBucket) Do(rw http.ResponseWriter, req *http.Request) {
	for _ = range sb.serversAvailable {
		server := sb.getServer()
		if !server.available {
			continue
		}
		req.Host = server.Location.Host
		req.URL.Host = server.Location.Host
		req.URL.Scheme = server.Location.Scheme
		req.RequestURI = ""
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			_, _ = fmt.Fprintln(rw, err)
			return
		}
		if response.StatusCode == 200 {
			rw.WriteHeader(response.StatusCode)
			io.Copy(rw, response.Body)
			return
		}
		server.excludeWithTimeout()
	}
	rw.WriteHeader(500)
}

func NewServerBucket() (sb *WeightedServerBucket) {
	var servers []weightedServer
	var serversAvailable []*weightedServer
	serverFile, _ := os.OpenFile("server-list", os.O_RDONLY, 0666)
	reader := bufio.NewScanner(serverFile)
	for reader.Scan() {
		directives := strings.Split(reader.Text(), ` `)
		address := directives[0]
		weight, err := strconv.Atoi(directives[1])
		timeout, err := strconv.Atoi(directives[2])
		host, err := url.Parse(address)
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, weightedServer{host, true, time.Duration(timeout) * time.Second, weight})
	}

	for serverId := range servers {
		for i := 0; i < servers[serverId].weight; i++ {
			serversAvailable = append(serversAvailable, &servers[serverId])
		}
	}
	sb = &WeightedServerBucket{servers, 0, serversAvailable}
	return
}

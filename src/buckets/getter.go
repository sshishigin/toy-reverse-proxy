package buckets

import (
	"net/http"
	"toy-reverse-proxy/src/buckets/simple"
	"toy-reverse-proxy/src/buckets/weighted"
)

type ServerBucket interface {
	Do(http.ResponseWriter, *http.Request)
}

type Server interface {
	excludeWithTimeout()
}

func GetServerBucket(bucketType string) ServerBucket {
	var bucket ServerBucket
	switch bucketType {
	case "simple":
		bucket = simple.NewServerBucket()
	case "weighted":
		bucket = weighted.NewServerBucket()
	default:
		panic("Pass a proper type of bucket")
	}
	return bucket
}

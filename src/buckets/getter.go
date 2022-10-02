package buckets

import "net/http"

type ServerBucket interface {
	Do(http.ResponseWriter, *http.Request)
}

type Server interface {
	excludeWithTimeout()
}

func GetServerBucket(bucketType string) ServerBucket {
	var Bucket ServerBucket
	switch bucketType {
	case "simple":
		Bucket = newSimpleServerBucket()
	case "weighted":
		Bucket = newWeightedServerBucket()
	default:
		panic("Pass a proper type of bucket")
	}
	return Bucket
}

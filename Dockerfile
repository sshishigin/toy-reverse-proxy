FROM golang:1.18.4
WORKDIR /app
COPY . .
RUN go build toy-reverse-proxy/cmd/simple
RUN echo "http://tokenizer:8000 5 150 5\nhttp://tokenizer_replica:8000 5 150 5" > server-list
ENTRYPOINT ["./simple"]
example: build buf-example

buf:
	buf generate --path ./protoc-gen-go-nats

buf-example:
	buf generate --template ./buf.gen.example.yaml --path ./example

build: 
	go build protoc-gen-go-nats/main.go

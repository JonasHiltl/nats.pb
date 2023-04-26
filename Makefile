example: build buf-example

buf:
	buf generate --path ./protoc-gen-nats

buf-example:
	buf generate --template ./buf.gen.example.yaml --path ./example

build: 
	go build protoc-gen-nats/main.go

# NATS-pb

`natspb` is a plugin for the Google protocol buffers compiler [protoc](https://github.com/protocolbuffers/protobuf). It generates [NATS Microservices](https://github.com/nats-io/nats.go/tree/main/micro) and clients of protobuf service definitions.

Through the NATS request-reply pattern we have support for load balancing and service discovery out of the box.  
Combined with protobuf as the data serialization format, `natspb` provides a simple and efficient way to build microservices using NATS.

## Installation
```
go get github.com/jonashiltl/natspb
```

## Usage
First specify your protobuf service with the `subject` of `protoc_gen_nats.options.nats` for each method set.

```protobuf
syntax = "proto3";

package example;
option go_package = "github.com/jonashiltl/natspb/example";

import "protoc-gen-nats/options/descriptor.proto";

service ExampleService {
    option(nats.pb.protoc_gen_nats.options.nats_service) = {
    name: "Example";
    description: "I'm a useful description";
    version: "1.0.0"
  };

  rpc Echo(Hello) returns (Hello) {
    option (nats.pb.protoc_gen_nats.options.nats) = {
      subject: "echo.echo";
    };
  };
}

message Hello {
  string greeting = 1;
}
```
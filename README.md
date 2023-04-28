# NATS-pb

`natspb` is a plugin for the Google protocol buffers compiler [protoc](https://github.com/protocolbuffers/protobuf). It generates [NATS Microservices](https://github.com/nats-io/nats.go/tree/main/micro) and clients of protobuf service definitions.

Through the NATS request-reply pattern we have support for load balancing and service discovery out of the box.  
Combined with protobuf as the data serialization format, `natspb` provides a simple and efficient way to build microservices using NATS.
## Installation

### Compile from source
The following instructions assume you are using
[Go Modules](https://github.com/golang/go/wiki/Modules) for dependency
management. Use a
[tool dependency](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
to track the versions of the following executable packages:

```go
// +build tools

package tools

import (
    _ "github.com/jonashiltl/nats.pb/protoc-gen-go-nats"
)
```

Run `go mod tidy` to resolve the versions. Install by running

```sh
go install github.com/jonashiltl/nats.pb/protoc-gen-go-nats
```

This will place four binaries in your `$GOBIN`;

- `protoc-gen-go-nats`

Make sure that your `$GOBIN` is in your `$PATH`.

## Usage

### Buf
Add `protoc-gen-go-nats` to your `buf.gen.yaml`
```yaml
version: v1
plugins:
  - plugin: go
    out: .
    opt: paths=source_relative
  - plugin: go-nats
    out: .
    opt: paths=source_relative
```
### Protoc
```
protoc -I . --go-nats_out ./gen/go \
    your/service/v1/your_service.proto
```

First specify your protobuf service with the `subject` of `protoc_gen_nats.options.nats` for each method set.  
You can optionally set the `protoc_gen_nats.options.nats_service` options which will be used when registering your service on NATS.

```protobuf
syntax = "proto3";

package example;
option go_package = "github.com/jonashiltl/nats.pb/example";

import "protoc-gen-go-nats/options/descriptor.proto";

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
// Code generated by protoc-gen-nats. DO NOT EDIT.

package example

import (
	"context"
	"time"
	"log"

	nats "github.com/nats-io/nats.go"
	natspb "github.com/jonashiltl/nats.pb"
	micro "github.com/nats-io/nats.go/micro"
	proto "google.golang.org/protobuf/proto"
	prototext "google.golang.org/protobuf/encoding/prototext"
)

type ExampleServiceClient interface {
	Echo(ctx context.Context, in *Hello, timeout time.Duration) (*Hello, error)
}

type exampleServiceClient struct {
	tr natspb.Transport
}

func NewExampleServiceClient(tr natspb.Transport) ExampleServiceClient {
	return &exampleServiceClient{tr}
}

type invokeParams struct {
	ctx     context.Context
	subj    string
	in      proto.Message
	timeout time.Duration
	out     proto.Message
}

func (c *exampleServiceClient) invoke(params invokeParams) error {
	b, err := proto.Marshal(params.in)
	if err != nil {
		return err
	}

	msg, err := c.tr.Request(params.subj, b, params.timeout)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(msg.Data, params.out)
	if err != nil {
		return err
	}
	return nil
}

func (c *exampleServiceClient) Echo(ctx context.Context, in *Hello, timeout time.Duration) (*Hello, error) {
	out := new(Hello)
	params := invokeParams{
		ctx:     ctx,
		subj:    "echo.echo",
		in:      in,
		timeout: timeout,
		out:     out,
	}
	err := c.invoke(params)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ExampleServiceServer interface {
	Echo(ctx context.Context, in *Hello) (*Hello, error)
}

func RegisterExampleServiceServer(nc *nats.Conn, srv ExampleServiceServer) (micro.Service, error) {
	s, err := micro.AddService(nc, micro.Config{
		Name:        "Example",
		Description: "I'm a useful description",
		Version:     "1.0.0",
	})
	if err != nil {
		return nil, err
	}

	// TODO: decide how to allow passing in context
	err = s.AddEndpoint(
		"Echo",
		micro.ContextHandler(context.Background(), _ExampleService_Echo_Handler(srv.Echo)),
		micro.WithEndpointSchema(&micro.Schema{
			Request:  prototext.Format(new(Hello)),
			Response: prototext.Format(new(Hello)),
		}),
	)
	if err != nil {
		log.Println(err)
	}

	return s, nil
}

func _ExampleService_Echo_Handler(mth func(context.Context, *Hello) (*Hello, error)) func(context.Context, micro.Request) {
	return func(ctx context.Context, r micro.Request) {
		in := new(Hello)
		err := proto.Unmarshal(r.Data(), in)

		msg, err := mth(ctx, in)

		res, err := proto.Marshal(msg)
		err = r.Respond(res)
	}
}

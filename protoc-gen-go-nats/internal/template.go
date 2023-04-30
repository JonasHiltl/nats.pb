package internal

import (
	"bytes"
	"log"
	"text/template"

	"github.com/jonashiltl/nats.pb/protoc-gen-go-nats/internal/utils"
	"github.com/jonashiltl/nats.pb/protoc-gen-go-nats/options"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var (
	clientTemplate = template.Must(template.New("client-template").Funcs(template.FuncMap{
		"toFirstLowerCase": utils.ToFirstLowerCase,
	}).Parse(`
	type {{ .GoServiceName }}Client interface {
		{{ range .Handlers -}}
			{{ .Name }} (ctx context.Context, in *{{ .RequestName }}, timeout time.Duration) (*{{ .ResponseName }}, error)
		{{ end }}
	}

	type {{ toFirstLowerCase .GoServiceName }}Client struct {
		tr natspb.Transport
	}

	func New{{ .GoServiceName }}Client(tr natspb.Transport) {{ .GoServiceName }}Client {
		return &{{ toFirstLowerCase .GoServiceName }}Client{tr}
	}

	type invokeParams struct {
		ctx context.Context;
		subj string; 
		in proto.Message;
		timeout time.Duration; 
		out proto.Message
	}
	
	func (c *{{ toFirstLowerCase .GoServiceName }}Client) invoke(params invokeParams) error {
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

	{{ range .Handlers -}}
		func (c *{{ toFirstLowerCase $.GoServiceName }}Client) {{ .Name }}(ctx context.Context, in *{{ .RequestName }}, timeout time.Duration) (*{{ .ResponseName }}, error) {
			out := new({{ .ResponseName }})
			params := invokeParams{
				ctx: ctx,
				subj: "{{ .Subject }}",
				in: in,
				timeout: timeout,
				out: out,
			}
			err := c.invoke(params)
			if err != nil {
				return nil, err
			}
			return out, nil
		}
	{{ end }}
	`))

	serverTemplate = template.Must(template.New("server-template").Parse(`
	type {{ .GoServiceName }}Server interface {
		{{ range .Handlers -}}
			{{ .Name }} (ctx context.Context, in *{{ .RequestName }}) (*{{ .ResponseName }}, error)
		{{ end }}
	}

	func Register{{ .GoServiceName }}Server(nc *nats.Conn, srv {{ .GoServiceName }}Server) (micro.Service, error) {
		s, err := micro.AddService(nc, micro.Config{
			Name: "{{ .MicroServiceName }}",
			Description: "{{ .MicroServiceDesc }}",
			Version: "{{ .MicroServiceVersion }}",
		})
		if err != nil {
			return nil, err
		}

		{{ range .Handlers -}}
			// TODO: decide how to allow passing in context
			err = s.AddEndpoint(
				"{{ .Name }}",
				micro.ContextHandler(context.Background(),  _{{ $.GoServiceName }}_{{ .Name }}_Handler(srv.{{ .Name }})),
				micro.WithEndpointSchema(&micro.Schema{
					Request: prototext.Format(new({{ .RequestName }})),
					Response: prototext.Format(new({{ .ResponseName }})),
				}),
			)
			if err != nil {
				log.Println(err)
			}
		{{ end }}


		return s, nil
	}

	{{ range .Handlers -}}
		func _{{ $.GoServiceName }}_{{ .Name }}_Handler(mth func(context.Context, *{{ .RequestName }}) (*{{ .ResponseName }}, error)) func(context.Context, micro.Request) {
			return func(ctx context.Context, r micro.Request) {
				in := new({{.RequestName}})
				_ = proto.Unmarshal(r.Data(), in)

				msg, _ := mth(ctx, in)

				res, _ := proto.Marshal(msg)
				_ = r.Respond(res)
			}
		}
	{{ end }}

	`))
)

type Handler struct {
	Name         string
	RequestName  string
	ResponseName string
	Subject      string
}

type clientTemplateParams struct {
	GoServiceName string
	Handlers      []Handler
}

type serverTemplateParams struct {
	GoServiceName       string
	MicroServiceName    string
	MicroServiceDesc    string
	MicroServiceVersion string
	Handlers            []Handler
}

type applyTemplateParams struct {
	srv *protogen.Service
}

func applyTemplate(params applyTemplateParams) (string, error) {
	w := bytes.NewBuffer(nil)

	handlers := make([]Handler, len(params.srv.Methods))
	for i, mth := range params.srv.Methods {
		opts := getMethodOptions(mth)
		subj := opts.Subject
		if subj == "" {
			log.Fatalln("Subject property of 'nats' option must be specified.")
		}

		handlers[i] = Handler{
			Name:         mth.GoName,
			RequestName:  mth.Input.GoIdent.GoName,
			ResponseName: mth.Output.GoIdent.GoName,
			Subject:      subj,
		}
	}

	cp := clientTemplateParams{
		GoServiceName: params.srv.GoName,
		Handlers:      handlers,
	}
	if err := clientTemplate.Execute(w, cp); err != nil {
		return "", err
	}

	opts := getServiceOptions(params.srv)
	sp := serverTemplateParams{
		GoServiceName: params.srv.GoName,
		// default micro service name to it's go name
		MicroServiceName: params.srv.GoName,
		Handlers:         handlers,
	}
	if opts != nil {
		if opts.Name != nil {
			sp.MicroServiceName = *opts.Name
		}
		if opts.Description != nil {
			sp.MicroServiceDesc = *opts.Description
		}
		if opts.Version != nil {
			sp.MicroServiceVersion = *opts.Version
		}
	}
	if err := serverTemplate.Execute(w, sp); err != nil {
		return "", err
	}

	return w.String(), nil
}

func getMethodOptions(mth *protogen.Method) *options.NatsMethodOptions {
	ext := proto.GetExtension(mth.Desc.Options(), options.E_Nats)
	opts, ok := ext.(*options.NatsMethodOptions)
	if !ok || opts == nil {
		log.Fatalln("Method option 'nats' is missing.")
	}

	return opts
}

func getServiceOptions(srv *protogen.Service) *options.NatsServiceOptions {
	ext := proto.GetExtension(srv.Desc.Options(), options.E_NatsService)
	opts, ok := ext.(*options.NatsServiceOptions)
	if !ok || opts == nil {
		log.Fatalln("Service option 'nats_service' is missing.")
	}
	return opts
}

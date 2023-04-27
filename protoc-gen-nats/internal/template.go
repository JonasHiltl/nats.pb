package internal

import (
	"bytes"
	"text/template"

	"github.com/jonashiltl/proto-nats/protoc-gen-nats/internal/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	clientTemplate = template.Must(template.New("client-template").Funcs(template.FuncMap{
		"toFirstLowerCase": utils.ToFirstLowerCase,
	}).Parse(`
	type {{ .ServiceName }}Client interface {
		{{ range .Endpoints -}}
			{{ .Name }} (ctx context.Context, in *{{ .InputName }}, timeout time.Duration) (*{{ .OutputName }}, error)
		{{ end }}
	}

	type {{ toFirstLowerCase .ServiceName }}Client struct {
		tr transport.Transport
	}

	func New{{ .ServiceName }}Client(tr transport.Transport) {{ .ServiceName }}Client {
		return &{{ toFirstLowerCase .ServiceName }}Client{tr}
	}

	
	type invokeParams struct {
		ctx context.Context;
		subj string; 
		in proto.Message;
		timeout time.Duration; 
		out proto.Message
	}
	
	func (c *{{ toFirstLowerCase $.ServiceName }}Client) invoke(params invokeParams) error {
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

	{{ range .Endpoints -}}
		func (c *{{ toFirstLowerCase $.ServiceName }}Client) {{ .Name }}(ctx context.Context, in *{{ .InputName }}, timeout time.Duration) (*{{ .OutputName }}, error) {
			out := new({{ .InputName }})
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
)

type Endpoint struct {
	Name       string
	InputName  string
	OutputName string
	Subject    string
}

type clientTemplateParams struct {
	ServiceName string
	Endpoints   []Endpoint
}

type applyTemplateParams struct {
	srv *protogen.Service
}

func applyTemplate(params applyTemplateParams) (string, error) {
	w := bytes.NewBuffer(nil)

	eps := make([]Endpoint, len(params.srv.Methods))
	for i, ep := range params.srv.Methods {
		// TODO: figure our how to read custom plugin options
		eps[i] = Endpoint{Name: ep.GoName, InputName: ep.Input.GoIdent.GoName, OutputName: ep.Output.GoIdent.GoName}
	}

	cp := clientTemplateParams{
		ServiceName: params.srv.GoName,
		Endpoints:   eps,
	}

	if err := clientTemplate.Execute(w, cp); err != nil {
		return "", err
	}

	return w.String(), nil
}

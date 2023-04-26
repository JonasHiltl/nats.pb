package internal

import (
	"bytes"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	clientTemplate = template.Must(template.New("client-template").Parse(`
	type {{ .ServiceName }}Client interface {
		{{ range .Endpoints -}}
			{{ .Name }} (ctx context.Context, in *{{ .InputName }}) (*{{ .OutputName }}, error)
		{{ end }}
	}
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

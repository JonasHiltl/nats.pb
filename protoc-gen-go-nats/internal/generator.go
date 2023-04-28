package internal

import (
	"go/format"

	"google.golang.org/protobuf/compiler/protogen"
)

type generator struct {
}

func New() *generator {
	return &generator{}
}

// generateFile generates a _nats.pb.go file containing nats service definitions.
func (g *generator) GenerateFile(gen *protogen.Plugin, file *protogen.File) (*protogen.GeneratedFile, error) {
	filename := file.GeneratedFilenamePrefix + "_nats.pb.go"
	gf := gen.NewGeneratedFile(filename, file.GoImportPath)
	gf.P("// Code generated by protoc-gen-nats. DO NOT EDIT.")
	gf.P()
	gf.P("package ", file.GoPackageName)
	gf.P()

	gf.P(`
		import (
			"context"
			"time"
			"log"
			
			nats "github.com/nats-io/nats.go"
			natspb "github.com/jonashiltl/nats.pb"
			micro "github.com/nats-io/nats.go/micro"		
			proto "google.golang.org/protobuf/proto"
		)
	`)

	for _, srv := range file.Services {
		code, err := applyTemplate(applyTemplateParams{
			srv: srv,
		})
		if err != nil {
			return nil, err
		}

		formatted, err := format.Source([]byte(code))
		if err != nil {
			return nil, err
		}

		_, err = gf.Write(formatted)
		if err != nil {
			return nil, err
		}

	}

	return gf, nil
}
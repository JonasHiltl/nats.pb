package main

import (
	"flag"
	"log"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/jonashiltl/natspb/protoc-gen-nats/internal"
)

func main() {
	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			g := internal.New()

			_, err := g.GenerateFile(gen, f)
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
}

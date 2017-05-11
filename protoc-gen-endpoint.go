package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/sermodigital/protoc-gen-endpoint/tables"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	eproto "github.com/sermodigital/protoc-gen-endpoint/proto"
	options "google.golang.org/genproto/googleapis/api/annotations"
)

func main() {
	log.SetFlags(255)

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	req := new(plugin.CodeGeneratorRequest)
	err = proto.Unmarshal(b, req)
	if err != nil {
		log.Fatalln(err)
	}

	err = writeEndpoints(os.Stdout, req)
	if err != nil {
		log.Fatalln(err)
	}
}

func writeEndpoints(w io.Writer, req *plugin.CodeGeneratorRequest) error {
	info, err := getInfo(req)
	if err != nil {
		return err
	}

	t, err := template.New("tmpl").Parse(templ)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, info)
	if err != nil {
		return err
	}

	fname := fmt.Sprintf("%s/%s.pb.ep.go", info.PkgName, info.PkgName)
	b, err := proto.Marshal(&plugin.CodeGeneratorResponse{
		File: []*plugin.CodeGeneratorResponse_File{
			{Name: proto.String(fname), Content: proto.String(buf.String())},
		},
	})
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

type Info struct {
	PkgName string
	Table   tables.Table
}

func getInfo(req *plugin.CodeGeneratorRequest) (inf Info, err error) {
	// From CodeGeneratorRequest's documentation:
	//
	// "FileDescriptorProtos for all files in files_to_generate and everything
	// 	they import.  The files will appear in topological order, so each file
	// 	appears before any file that imports it."
	pfs := req.GetProtoFile()
	pf := pfs[len(pfs)-1]

	inf.PkgName = pkgName(pf)
	inf.Table = make(tables.Table)

	for _, srv := range pf.GetService() {
		for _, meth := range srv.GetMethod() {
			if meth.Options == nil ||
				!proto.HasExtension(meth.Options, options.E_Http) {
				continue
			}

			ext, err := proto.GetExtension(meth.Options, options.E_Http)
			if err != nil {
				return inf, err
			}
			http, ok := ext.(*options.HttpRule)
			if !ok {
				return inf, fmt.Errorf("got %T, wanted *options.HttpRule", ext)
			}

			ext, _ = proto.GetExtension(meth.Options, eproto.E_Endpoint)
			endp, ok := ext.(*eproto.Endpoint)
			unauth := ok && endp.Unauthenticated

			prefix := strings.TrimSuffix(inf.PkgName, "pb")
			action := prefix + "." + *meth.Name

			err = parseTuple(http, inf.Table, unauth, action)
			if err != nil {
				return inf, err
			}

			for _, http := range http.AdditionalBindings {
				err := parseTuple(http, inf.Table, unauth, action)
				if err != nil {
					return inf, err
				}
			}
		}
	}
	return inf, nil
}

// pkgName returns a suitable package name from file.
//
// Mostly borrowed from grpc-gateway.
func pkgName(file *descriptor.FileDescriptorProto) string {
	if file.Options != nil && file.Options.GoPackage != nil {
		gopkg := file.Options.GetGoPackage()
		i := strings.LastIndexByte(gopkg, '/')
		if i < 0 {
			return gopkg
		}
		return strings.Replace(gopkg[i+1:], ".", "_", -1)
	}
	if file.Package == nil {
		base := filepath.Base(file.GetName())
		ext := filepath.Ext(base)
		return strings.TrimSuffix(base, ext)
	}
	return strings.Replace(file.GetPackage(), ".", "_", -1)
}

// parseTuple parses a new tables.Endpoint from http and adds it to table.
func parseTuple(http *options.HttpRule, tbl tables.Table, unauth bool, action string) error {
	var url string
	act := tables.Action{Unauthenticated: unauth, Name: action}
	switch v := http.Pattern.(type) {
	case *options.HttpRule_Get:
		url = v.Get
		act.Method = "GET"
	case *options.HttpRule_Put:
		url = v.Put
		act.Method = "PUT"
	case *options.HttpRule_Post:
		url = v.Post
		act.Method = "POST"
	case *options.HttpRule_Delete:
		url = v.Delete
		act.Method = "DELETE"
	case *options.HttpRule_Patch:
		url = v.Patch
		act.Method = "PATCH"
	case *options.HttpRule_Custom:
		url = v.Custom.Path
		act.Method = v.Custom.Kind
	default:
		return fmt.Errorf("unknown http.Pattern: %T", http.Pattern)
	}
	ep := tbl[url]
	ep.Add(act)
	tbl[url] = ep
	return nil
}

const templ = `// Package {{ .PkgName }} creates a (URL, HTTP method) -> action lookup table
package {{ .PkgName }}

import "github.com/sermodigital/protoc-gen-endpoint/tables"

// Table returns a tables.Table containing the endpoints within a gRPC package.
func Table() tables.Table {
	return tables.Table{
		{{- range $url, $eps := .Table }}
		{{ $url | printf "%q" }}: tables.Endpoint{
			Methods: {{ $eps.Methods | printf "%q" -}},
			Actions: []tables.Action{
				{{ range $act := $eps.Actions -}}
				{
					Name: {{ $act.Name | printf "%q" -}},
					Method: {{- $act.Method | printf "%q" }},
					Unauthenticated: {{ $act.Unauthenticated }},
				},
				{{- end }}
			},
		},
		{{- end  }}
	}
}
`

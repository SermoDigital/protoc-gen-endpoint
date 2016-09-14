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

	"github.com/SermoDigital/protoc-gen-endpoint/tables"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	eproto "github.com/SermoDigital/protoc-gen-endpoint/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	options "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"
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
	infos, err := getInfo(req)
	if err != nil {
		return err
	}
	t, err := template.New("tmpl").Parse(templ)
	if err != nil {
		return err
	}

	var files []*plugin.CodeGeneratorResponse_File
	var buf bytes.Buffer
	for _, info := range infos {
		buf.Reset()
		err = t.Execute(&buf, info)
		if err != nil {
			return err
		}
		fname := fmt.Sprintf("%s/%s.pb.ep.go", info.PkgName, info.PkgName)
		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fname),
			Content: proto.String(buf.String()),
		})
	}

	b, err := proto.Marshal(&plugin.CodeGeneratorResponse{File: files})
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

func getInfo(req *plugin.CodeGeneratorRequest) (ifs []Info, err error) {
	for _, pf := range req.GetProtoFile() {
		i := Info{
			PkgName: pkgName(pf),
			Table:   make(tables.Table),
		}
		for _, srv := range pf.GetService() {
			for _, meth := range srv.GetMethod() {
				if meth.Options == nil ||
					!proto.HasExtension(meth.Options, options.E_Http) ||
					!proto.HasExtension(meth.Options, eproto.E_Endpoint) {
					continue
				}

				ext, err := proto.GetExtension(meth.Options, options.E_Http)
				if err != nil {
					return nil, err
				}
				http, ok := ext.(*options.HttpRule)
				if !ok {
					return nil, fmt.Errorf("got %T, wanted *options.HttpRule", ext)
				}

				ext, err = proto.GetExtension(meth.Options, eproto.E_Endpoint)
				if err != nil {
					return nil, err
				}
				endp, ok := ext.(*eproto.Endpoint)
				if !ok {
					return nil, fmt.Errorf("got %T, wanted *eproto.Endpoint", ext)
				}

				err = parseTuple(http, i.Table, endp.Unauthenticated, "")
				if err != nil {
					return nil, err
				}

				for _, http := range http.AdditionalBindings {
					err := parseTuple(http, i.Table, endp.Unauthenticated, "")
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if len(i.Table) != 0 {
			ifs = append(ifs, i)
		}
	}
	return ifs, nil
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
	var (
		url    string
		method string
	)
	switch v := http.Pattern.(type) {
	case *options.HttpRule_Get:
		url = v.Get
		method = "GET"
	case *options.HttpRule_Put:
		url = v.Put
		method = "PUT"
	case *options.HttpRule_Post:
		url = v.Post
		method = "POST"
	case *options.HttpRule_Delete:
		url = v.Delete
		method = "DELETE"
	case *options.HttpRule_Patch:
		url = v.Patch
		method = "PATCH"
	case *options.HttpRule_Custom:
		url = v.Custom.Path
		method = v.Custom.Kind
	default:
		return fmt.Errorf("unknown http.Patten: %T", http.Pattern)
	}

	eps := tbl[url]
	eps = append(eps, tables.Endpoint{
		Method:          method,
		Unauthenticated: unauth,
		Action:          action,
	})
	tbl[url] = eps
	return nil
}

const templ = `// Package {{ .PkgName }} creates a (URL, HTTP method) -> action lookup table
package {{ .PkgName }}

import "github.com/SermoDigital/protoc-gen-endpoint/tables"

// Table returns a tables.Table containing the endpoints within a gRPC package.
func Table() tables.Table {
	return tables.Table{
		{{ range $url, $eps := .Table -}}
		{{- $url | printf "%q" }}: []tables.Endpoint{
			{{ range $ep := $eps -}}
			{
				Method: {{- $ep.Method | printf "%q" }},
				Unauthenticated: {{ $ep.Unauthenticated }},
				Action: {{ $ep.Action | printf "%q" -}},
			},
			{{- end }}
		},
		{{- end  }}
	}
}
`

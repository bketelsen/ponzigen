package main

import (
	"bytes"
	"go/types"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ernesto-jimenez/gogen/cleanimports"
	"github.com/ernesto-jimenez/gogen/gogenutil"
	"github.com/ernesto-jimenez/gogen/importer"
)

// Generator will generate the Ponzi code
type Generator struct {
	name    string
	namePkg string
	pkgPath string
	pkg     *types.Package
	targets map[string]*types.Struct
}

// NewGenerator initializes a Generator
func NewGenerator(pkg string) (*Generator, error) {
	var err error
	if pkg == "" || pkg[0] == '.' {
		pkg, err = filepath.Abs(filepath.Clean(pkg))
		if err != nil {
			return nil, err
		}
		pkg = gogenutil.StripGopath(pkg)
	}
	p, err := importer.Default().Import(pkg)
	if err != nil {
		return nil, err
	}
	targets := make(map[string]*types.Struct)
	for _, name := range p.Scope().Names() {

		obj := p.Scope().Lookup(name)
		if _, ok := obj.Type().Underlying().(*types.Struct); ok {
			if !strings.Contains(name, "ListResult") {
				targets[name] = obj.Type().Underlying().(*types.Struct)
			}
		}

	}
	return &Generator{
		pkg:     p,
		targets: targets,
		pkgPath: pkg,
	}, nil
}

func (g Generator) qf(pkg *types.Package) string {
	if g.pkg == pkg {
		return ""
	}
	return pkg.Name()
}

func (g Generator) SourcePackage() string {
	return g.pkg.Name()
}

func (g Generator) Package() string {
	if g.namePkg != "" {
		return g.namePkg
	}
	return g.pkg.Name()
}

func (g Generator) PackagePath() string {
	return g.pkgPath
}

func (g *Generator) SetPackage(name string) {
	g.namePkg = name
}
func (g Generator) Targets() map[string]*types.Struct {
	return g.targets
}
func (g Generator) Write(wr io.Writer) error {
	var buf bytes.Buffer
	if err := ponziTmpl.Execute(&buf, g); err != nil {
		return err
	}
	return cleanimports.Clean(wr, buf.Bytes())
}

var (
	funcMap = template.FuncMap{
		"lower": strings.ToLower,
	}
	ponziTmpl = template.Must(template.New("ponzi").Funcs(funcMap).Parse(`/*
* CODE GENERATED AUTOMATICALLY WITH github.com/bketelsen/ponzigen
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/
{{$sp := .SourcePackage}}
package {{.Package}}

import (
	"github.com/bketelsen/ponzi"
	"time"
	"github.com/pkg/errors"
	"{{.PackagePath}}"
)

var BaseURL string 

{{ range $n,$s := .Targets }}type {{$n}}ListResult struct {
	Data []{{$sp}}.{{$n}} ` + "`" + "json:" + "\"" + "data" + "\"`" + `
}
{{ end }}

{{ range $n,$s := .Targets }}var {{lower $n}}Cache *ponzi.Cache
{{ end }}

{{ range $n,$s := .Targets }}func init{{$n}}Cache() {
	if {{lower $n}}Cache == nil {
		{{lower $n}}Cache = ponzi.New(BaseURL, 1*time.Minute, 30*time.Second)
	}
}
{{ end }}

{{ range $n,$s := .Targets }}func Get{{$n}}(id int) ({{$sp}}.{{$n}}, error) {
	init{{$n}}Cache()
	var sp {{$n}}ListResult
	err := {{lower $n}}Cache.Get(id, "{{$n}}", &sp)
	if err != nil {
		return {{$sp}}.{{$n}}{}, err
	}
	if len(sp.Data) == 0 {
		return {{$sp}}.{{$n}}{}, errors.New("Not Found")
	}
	return sp.Data[0], err

}
{{ end }}
{{ range $n,$s := .Targets }}func Get{{$n}}BySlug(slug string) ({{$sp}}.{{$n}}, error) {
	init{{$n}}Cache()
	var sp {{$n}}ListResult
	err := {{lower $n}}Cache.GetBySlug(slug, "{{$n}}", &sp)
	if err != nil {
		return {{$sp}}.{{$n}}{}, err
	}
	if len(sp.Data) == 0 {
		return {{$sp}}.{{$n}}{}, errors.New("Not Found")
	}
	return sp.Data[0], err

}
{{ end }}
{{ range $n,$s := .Targets }}func Get{{$n}}BySlug(slug string) ({{$sp}}.{{$n}}, error) {
       init{{$n}}Cache()
       var sp {{$n}}ListResult
       err := {{lower $n}}Cache.GetBySlug(slug, "{{$n}}", &sp)
       if err != nil {
               return {{$sp}}.{{$n}}{}, err
       }
       if len(sp.Data) == 0 {
              return {{$sp}}.{{$n}}{}, errors.New("Not Found")
       }
       return sp.Data[0], err

}
{{ end }}
{{ range $n,$s := .Targets }}func Get{{$n}}List() ([]{{$sp}}.{{$n}}, error) {
	init{{$n}}Cache()
	var sp {{$n}}ListResult
	err := {{lower $n}}Cache.GetAll("{{$n}}", &sp)
	if err != nil {
		return []{{$sp}}.{{$n}}{}, err
	}
	if len(sp.Data) == 0 {
		return []{{$sp}}.{{$n}}{}, errors.New("Not Found")
	}
	return sp.Data, err

}
{{ end }}
`))
)

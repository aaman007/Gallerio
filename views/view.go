package views

import (
	"bytes"
	"errors"
	"gallerio/utils/context"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   = "views/layouts/"
	TemplateDir = "views/"
	TemplateExt = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented properly")
		},
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	v.Render(w, req, nil)
}

func (v *View) Render(w http.ResponseWriter, req *http.Request, data interface{}) {
	var _data Data
	switch d := data.(type) {
	case Data:
		_data = d
	default:
		_data = Data{
			Content: data,
		}
	}

	if alert := getAlert(req); alert != nil {
		_data.Alert = alert
		clearAlert(w)
	}
	_data.User = context.User(req.Context())
	csrfField := csrf.TemplateField(req)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	
	var buff bytes.Buffer
	if err := tpl.ExecuteTemplate(&buff, v.Layout, _data); err != nil {
		log.Println(err)
		http.Error(w, AlertMessageGeneric, http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buff)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func addTemplatePath(files []string) {
	for i, file := range files {
		files[i] = TemplateDir + file
	}
}

func addTemplateExt(files []string) {
	for i, file := range files {
		files[i] = file + TemplateExt
	}
}

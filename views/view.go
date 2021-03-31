package views

import (
	"bytes"
	"gallerio/utils/context"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir = "views/layouts/"
	TemplateDir = "views/"
	TemplateExt = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout: layout,
	}
}

type View struct {
	Template *template.Template
	Layout string
}

func (v *View) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	v.Render(w, req,nil)
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

	_data.User = context.User(req.Context())

	var buff bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buff, v.Layout, _data); err != nil {
		http.Error(w, AlertMessageGeneric, http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buff)
}

func layoutFiles() []string  {
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

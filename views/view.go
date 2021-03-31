package views

import (
	"bytes"
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
	v.Render(w, nil)
}

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	switch data.(type) {
	case Data:
		// pass
	default:
		data = Data{
			Content: data,
		}
	}
	var buff bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buff, v.Layout, data); err != nil {
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

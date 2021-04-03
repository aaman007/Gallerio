package main

import (
	"flag"
	"os"
	"text/template"
)

type data struct {
	Name   string
	Object string
}

func main() {
	// To Generate a model write the following command
	// run : go run cmd/gen/main.go -name=GeneratedModel -object=generatedModel  > models/qenerated_model.go
	
	var d data
	flag.StringVar(&d.Name, "name", "", "Name of the model in pascal case")
	flag.StringVar(&d.Object, "object", "", "Name of the object in camel case")
	flag.Parse()
	
	t := template.Must(template.New("model").Parse(modelTemplate))
	t.Execute(os.Stdout, d)
}

const modelTemplate = `
package models

import (
	"github.com/jinzhu/gorm"
)

type {{.Name}} struct {
	gorm.Model
}

type {{.Name}}DB interface {
	Create({{.Object}} *{{.Name}}) error
	Delete(id uint) error
}

func New{{.Name}}Service(db *gorm.DB) {{.Name}}Service {
	return &{{.Object}}Service{
		{{.Name}}DB: &{{.Object}}Validator{&{{.Object}}Gorm{db}},
	}
}

type {{.Name}}Service interface {
	{{.Name}}DB
}

type {{.Object}}Service struct {
	{{.Name}}DB
}

type {{.Object}}ValFunc func({{.Object}} *{{.Name}}) error

func run{{.Name}}ValFuncs({{.Object}} *{{.Name}}, fns ...{{.Object}}ValFunc) error {
	for _, fn := range fns {
		if err := fn({{.Object}}); err != nil {
			return err
		}
	}
	return nil
}

type {{.Object}}Validator struct {
	{{.Name}}DB
}

func (mv *{{.Object}}Validator) Create({{.Object}} *{{.Name}}) error {
	err := run{{.Name}}ValFuncs({{.Object}})
	if err != nil {
		return err
	}
	return mv.{{.Name}}DB.Create({{.Object}})
}

func (mv *{{.Object}}Validator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return mv.{{.Name}}DB.Delete(id)
}

type {{.Object}}Gorm struct {
	db *gorm.DB
}

func (mg *{{.Object}}Gorm) Create({{.Object}} *{{.Name}}) error {
	return mg.db.Create({{.Object}}).Error
}

func (mg *{{.Object}}Gorm) Delete(id uint) error {
	{{.Object}} := {{.Name}}{Model: gorm.Model{ID: id}}
	return mg.db.Delete(&{{.Object}}).Error
}
`

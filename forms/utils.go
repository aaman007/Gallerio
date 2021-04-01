package forms

import (
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
)

func ParseForm(req *http.Request, form interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	return ParseValues(req.PostForm, form)
}

func ParseURLParams(req *http.Request, form interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	return ParseValues(req.Form, form)
}

func ParseValues(values url.Values, form interface{}) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(form, values); err != nil {
		return err
	}
	return nil
}

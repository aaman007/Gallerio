package forms

import (
	"github.com/gorilla/schema"
	"net/http"
)

func ParseForm(req *http.Request, form interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	if err := decoder.Decode(form, req.PostForm); err != nil {
		return err
	}
	return nil
}

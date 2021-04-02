package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gallerio/models"
	"gallerio/utils/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

func NewOAuthsController(os models.OAuthService, configs map[string]*oauth2.Config) *OAuthsController {
	return &OAuthsController{
		os:      os,
		configs: configs,
	}
}

type OAuthsController struct {
	os      models.OAuthService
	configs map[string]*oauth2.Config
}

func (oc *OAuthsController) Connect(w http.ResponseWriter, req *http.Request) {
	provider := mux.Vars(req)["provider"]
	state := csrf.Token(req)
	if _, ok := oc.configs[provider]; !ok {
		http.Error(w, "Unknown Provider", http.StatusBadRequest)
		return
	}
	
	cookie := &http.Cookie{
		Name: "oauth_state",
		Value: state,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	
	url := oc.configs[provider].AuthCodeURL(state)
	http.Redirect(w, req, url, http.StatusFound)
}

func (oc *OAuthsController) Callback(w http.ResponseWriter, req *http.Request) {
	provider := mux.Vars(req)["provider"]
	if _, ok := oc.configs[provider]; !ok {
		http.Error(w, "Unknown Provider", http.StatusBadRequest)
		return
	}
	
	req.ParseForm()
	state := req.FormValue("state")
	cookie, err := req.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid State", http.StatusBadRequest)
		return
	}
	
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	
	code := req.FormValue("code")
	token, err := oc.configs[provider].Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	user := context.User(req.Context())
	existing, err := oc.os.Find(user.ID, provider)
	
	switch err {
	case models.ErrNotFound:
		// pass
	case nil:
		oc.os.Delete(existing.ID)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	oauth := &models.OAuth{
		UserID: user.ID,
		Provider: provider,
		Token: *token,
	}
	err = oc.os.Create(oauth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	fmt.Fprintln(w, token)
}

func (oc *OAuthsController) DropboxTest(w http.ResponseWriter, req *http.Request) {
	provider := mux.Vars(req)["provider"]
	if provider != models.OAuthDropbox {
		http.Error(w, "Unknown Provider", http.StatusBadRequest)
		return
	}
	
	req.ParseForm()
	path := req.FormValue("path")
	
	user := context.User(req.Context())
	oauth, err := oc.os.Find(user.ID, models.OAuthDropbox)
	if err != nil {
		panic(err)
	}
	
	token := oauth.Token
	data := struct {
		Path string `json:"path"`
	}{Path: path}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	
	client := oc.configs["provider"].Client(context.TODO(), &token)
	reqObj, err := http.NewRequest(
		http.MethodPost,
		"https://api.dropboxapi.com/2/files/list_folder",
		bytes.NewReader(dataBytes),
	)
	if err != nil {
		panic(err)
	}
	reqObj.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(reqObj)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

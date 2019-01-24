package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	llctx "lenslockedbr.com/context"
	"lenslockedbr.com/models"
)

type OAuths struct {
	os models.OAuthService
	configs map[string]*oauth2.Config
}

func NewOAuths(os models.OAuthService, configs map[string]*oauth2.Config) *OAuths {
	return &OAuths {
		os: os,
		configs: configs, 
	}
}

func(o *OAuths) Connect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	service := vars["service"]
	oauthCfg, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", 
                           http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	cookie := http.Cookie {
		Name: "oauth_state",
		Value: state,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	url := oauthCfg.AuthCodeURL(state)
	log.Println(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuths) Callback(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	service := vars["service"]
	oauthCfg, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", 
                           http.StatusBadRequest)
		return
	}

	r.ParseForm()

	fmt.Fprintln(w, "code: ", r.FormValue("code"),
                     " state: ", r.FormValue("state"))

	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state provided", 
                           http.StatusBadRequest)
		return
	}

	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	code := r.FormValue("code")
	token, err := oauthCfg.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := llctx.User(r.Context())
	existing, err := o.os.Find(user.ID, service)
	if err == models.ErrNotFound {
		// noop
	} else if err != nil {
		http.Error(w, err.Error(), 
                           http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(existing.ID)
	}

	userOAuth := models.OAuth {
		UserID: user.ID,
		Token: *token,
		Service: service,
	}	
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), 
                           http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%+v", token)
}


func (o *OAuths) DropboxTest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	service := vars["service"]
	oauthCfg, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", 
                           http.StatusBadRequest)
		return
	}

	r.ParseForm()
	path := r.FormValue("path")

	user := llctx.User(r.Context())
	userOAuth, err := o.os.Find(user.ID, service)
	if err != nil {
		panic(err)
	} 

	token := userOAuth.Token
	client := oauthCfg.Client(context.TODO(), &token)

	url := "https://api.dropboxapi.com/2/files/list_folder"

	data := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	} 

	req, err := http.NewRequest(http.MethodPost,
                                    url,
                                    bytes.NewReader(dataBytes))
	if err != nil {
		panic(err)
	} 

	req.Header.Add("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	} 

	defer resp.Body.Close()

	io.Copy(w, resp.Body)
	
}





package main

import (
	"amf_common"

	//"time"

	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"strings"
)

type appState struct{
	newstate string;
	oldstate string;
}

type Page struct {
	Title string
	Body  []byte
}

var appsStates map[string]appState

func main() {
	amf_common.AppHeader = amf_common.TAppHeader{
		AppName: "ph1_control",
		AppVersion: "1.0.0",
	}
	amf_common.OnAppStart()

	//load phonebook

	defer amf_common.OnAppEnd()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":" + amf_common.AppHeader.ControlPort, nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	lp,_:= loadPage(r.URL.Path[1:])
	//w.WriteHeader(http.StatusOK)
	if (lp == nil) {
		fmt.Fprintf(w, "")
		return
	}
	xt:= strings.Split(lp.Title, ".")
	ext:= xt[len(xt)-1]

	contentType := "text/html"

	if (ext == "css") {
		contentType = "text/css"
	}else if (ext == "jpg") {
		contentType = "image/jpeg"
	}else if (ext == "png") {
		contentType = "image/png"
	}else if (ext == "js") {
		contentType = "text/javascript"
	}


	w.Header().Set("Content-Type", contentType)
	fmt.Fprintf(w, "%s", lp.Body)
}

func loadPage(urlpath string) (*Page, error) {
	if ((urlpath == "") || (urlpath == "/")) {
		urlpath = "index.html"
	}
	filename := "pages/" + urlpath
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &Page{Title: urlpath, Body: body}, nil
}
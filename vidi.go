package vidi

import (
	"fmt"
	"net/http"
)

type VidiContext struct {
	systemTest string
}

func (v *VidiContext) ProcessBody() {
	fmt.Println("temp")
}

func (v *VidiContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/call" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Call Complete!")
}

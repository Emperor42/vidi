package vidi

import (
	"fmt"
	"net/http"
)

type VidiContext struct {
	Name string
}

func (v *VidiContext) ProcessBody() {
	fmt.Println("temp")
}

func (v *VidiContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/data" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		w.Header().Set("Content-Type", "text/html")
		value := []byte(v.Name + "\n")
		w.Write(value)
	}

	fmt.Fprintf(w, "Call Complete!")
	return
}

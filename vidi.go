package vidi

import (
	"fmt"
	"net/http"
)

type VidiInterface interface {
	Connect(any) any
	Create(any) any
	Read(any) any
	Update(any) any
	Delete(any) any
}

type VidiContext struct {
	Name string
}

func (v *VidiContext) ProcessBody() {
	fmt.Println("temp")
}

func (v *VidiContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() != "/data" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, v.Name)
	}

	fmt.Fprintf(w, "VIDI - Call Complete!")
}

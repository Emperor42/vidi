package vidi

import (
	"fmt"
	"log"
	"net/http"
)

type vidiContext struct {
	d             http.Dir
	generatedFile []string
	err           error
	veniTargets   []string
}

func (d *veniContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.generatedFile = []string{}
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)
	name := r.URL.Path
	file, err := d.Open(name)
	if err != nil {
		fmt.Println("Issue opening file, please review")
		log.Fatal(err)
	}
	err = d.processHTML(file)
	if err != nil {
		fmt.Println("Issue reading file, please review")
		log.Fatal(err)
	}
	//close my file
	fmt.Println("File processed!")
	//write line by line to resp
	charCount := 0
	d.generatedFile = []string{}
	w.Header().Set("Content-Type", "text/html")
	for i, v := range d.generatedFile {
		value := []byte(v + "\n")
		tmpCount, erWrite := w.Write(value)
		charCount += tmpCount
		if erWrite != nil {
			fmt.Printf("Issue writing please review at line: %n\n", i)
			log.Fatal(erWrite)
		} else {
			fmt.Printf("Wrote line: %n\n", i)
		}
	}
}

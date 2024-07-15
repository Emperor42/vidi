package vidi 

import(
	"net/http"
)

type Record struct {
	Head string 
	Body string
	Tail string
	Code string
	Mark string
	Sign string
	Line uint64
}

type Database struct {
	rows []Record
}

func Load(handler http.Handler) http.Handler{
	//right now do nothing
	return handler;
}
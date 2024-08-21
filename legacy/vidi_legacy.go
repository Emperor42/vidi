package vidi_legacy

import (
	"database/sql"
	"fmt"
	"net/http"
)

/*
*
Record
Head: the main source (key) of this data
Body: the main content (value) of this data
Tail: links to forms specific and general
Sign: Signature of the member who created the record
Mark: Datetime stamp with location where the source came from
Code: Request URL given and body
Line: Unique auto generated identifier for this line
*/
type Record struct {
	Head string
	Body string
	Tail string
	Code string
	Mark string
	Sign string
	Line uint64
}

/*
*
Access
Record: The record to which to grant access
Member: The member to whom to grant access
Number: Targeting record number value
*/
type Access struct {
	Record uint64
	Member uint64
	Number uint64
}

/*
*
Member
Name:
*/
type Member struct {
	Name string
	Line uint64
	Peer uint64
}

type Service struct {
	db  *sql.DB // Database connection pool.
	err error
}

func (s *Service) queryAllByHead(name string) ([]Record, error) {
	// An records slice to hold data from returned rows.
	var records []Record
	rows, err := s.db.Query("SELECT * FROM record WHERE recordhead = ?", name)
	if err != nil {
		return nil, fmt.Errorf("search %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Record
		if err := rows.Scan(&alb.Line, &alb.Sign, &alb.Code, &alb.Mark, &alb.Head, &alb.Body, &alb.Tail); err != nil {
			return nil, fmt.Errorf("search %q: %v", name, err)
		}
		records = append(records, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search %q: %v", name, err)
	}
	return records, nil
}

func (s *Service) query(name string) ([]Record, error) {
	// An records slice to hold data from returned rows.
	var records []Record
	rows, err := s.db.Query("SELECT * FROM record WHERE head = ?", name)
	if err != nil {
		return nil, fmt.Errorf("search %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Record
		if err := rows.Scan(&alb.Line, &alb.Sign, &alb.Code, &alb.Mark, &alb.Head, &alb.Body, &alb.Tail); err != nil {
			return nil, fmt.Errorf("search %q: %v", name, err)
		}
		records = append(records, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search %q: %v", name, err)
	}
	return records, nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//db := s.db
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)
	//name := r.URL.Path
	//write line by line to resp
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("VIDI"))
}

func Load(db *sql.DB) *Service {
	s := &Service{db: db, err: nil}
	return s
}

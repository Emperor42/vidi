VIDI Middleware for Go 🛡️

A lightweight HTTP middleware that enforces method restrictions, validates custom form data, and provides a simple JSON‑storage/pull mechanism.
Designed for the standard net/http package (works with any router that accepts http.Handler).
Table of Contents

    Features
    Installation
    Usage Example
    Middleware Behaviour
    Data Flow Details
    API Reference
    Extending the Middleware
    License

Features

    Method gating – only GET and POST are accepted; everything else returns 405.
    GET passthrough – GET requests flow straight to the next handler.
    POST validation – looks for a form field VIDI.
        Accepts only PUSH or PULL (case‑insensitive). Anything else yields 400.
    PUSH mode – converts all non‑VICI form fields to JSON, stores a mapping between the generated JSON “file” and the originating key/value pairs.
    PULL mode – finds the stored JSON record that matches the greatest number of incoming key/value pairs and returns it.
    Security – never persists a field named VICI (or its value).
    Thread‑safe in‑memory storage – safe for concurrent requests.

Installation

go get github.com/yourusername/vidi

(Replace github.com/yourusername/vidi with the actual module path where you place the source.)
Usage Example

package main

import (
	"log"
	"net/http"

	"github.com/yourusername/vidi" // import the middleware package
)

func main() {
	mux := http.NewServeMux()

	// Example endpoint that receives POSTs
	mux.HandleFunc("/api/endpoint", func(w http.ResponseWriter, r *http.Request) {
		// Business logic after VIDI processing
		w.Write([]byte("handler reached"))
	})

	// Wrap the router with the VIDI middleware
	handler := vidi.VIDI(mux)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

Run:

go run .

Now the server:

    Accepts GET requests unchanged.
    Handles POST requests according to the VIDI rules.

Middleware Behaviour
Step	Action	Outcome
1	Request arrives	Middleware checks HTTP method.
2	Method ≠ GET/POST	Returns 405 Method Not Allowed (with Allow: GET, POST).
3	Method = GET	Calls next.ServeHTTP – no further processing.
4	Method = POST	Parses form data (r.ParseForm()).
5	Missing VIDI key	Calls next.ServeHTTP – request proceeds untouched.
6	VIDI present	Value normalized to uppercase.
7	Value not PUSH/PULL	Returns 400 Bad Request (“invalid VIDI value”).
8	Strip any VICI field	Guarantees VICI never gets stored.
9a	VIDI = PUSH	- Convert remaining form fields to JSON.- Generate deterministic filename.- Store mapping in a thread‑safe map.- Respond with the JSON payload.
9b	VIDI = PULL	- Compare incoming fields against stored records.- Return the record with the highest number of matching key/value pairs (or 404 if none).
Data Flow Details
PUSH Workflow

    Clean form – remove VICI; keep first value of each remaining key.
    Marshal to JSON – pretty‑printed (json.MarshalIndent).
    Generate filename – simple checksum‑based name (vidi_<HEX>.json).
    Store – records[filename] = jsonRecord{FileName: filename, Content: cleanForm} (protected by sync.RWMutex).
    Response – Content-Type: application/json with the generated JSON.

PULL Workflow

    Clean form – same as PUSH.
    Iterate over stored records – count exact key/value matches (countMatchingPairs).
    Select best match – highest match count.
    Response – JSON object containing:

    {
      "file": "<filename>",
      "match_count": <int>,
      "total_stored": <int>,
      "data": { ...original key/value pairs... }
    }

API Reference
func VIDI(next http.Handler) http.Handler

Creates the middleware. Pass your router or handler as next.
func GetRecords() map[string]jsonRecord

Returns a copy of the internal map (filename → jsonRecord).
Useful for diagnostics or exposing the data via another endpoint.
Types

type jsonRecord struct {
    FileName string            // e.g. "vidi_AB12CD.json"
    Content  map[string]string // the original key/value pairs stored
}

Extending the Middleware
Goal	Suggested Change
Persist JSON to disk	Replace the in‑memory records map with a DB (SQLite, PostgreSQL) or write files using os.WriteFile.
Support multi‑value fields	Change map[string]string to map[string][]string and adjust countMatchingPairs.
Stronger filename hashing	Use crypto/sha256 and hex‑encode the digest instead of the simple checksum.
Custom error messages	Wrap http.Error calls with your own JSON error struct.
Rate limiting / auth	Insert additional middleware before/after VIDI.

All modifications stay within the same http.Handler contract, so they compose nicely with other Go middlewares.
License

MIT License – feel free to use, modify, and distribute. See the LICENSE file for details.

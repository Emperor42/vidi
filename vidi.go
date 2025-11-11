package vidi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
)

// -------------------------------------------------------------------
// 1️⃣  Middleware skeleton – reject everything except GET/POST
// -------------------------------------------------------------------
func VIDI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// --------------------------------------------------------------
		// 1. Allow only GET and POST, reject everything else
		// --------------------------------------------------------------
		switch r.Method {
		case http.MethodGet:
			// ----------------------------------------------------------
			// 2. Pass GET straight through – no extra work
			// ----------------------------------------------------------
			next.ServeHTTP(w, r)
			return

		case http.MethodPost:
			// ----------------------------------------------------------
			// 3. Process POST body (form data)
			// ----------------------------------------------------------
			if err := r.ParseForm(); err != nil {
				http.Error(w, "unable to parse form data", http.StatusBadRequest)
				return
			}
			handlePOST(w, r, next)

		default:
			// Anything else → 405 Method Not Allowed
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}

// -------------------------------------------------------------------
// Helper structures for requirement #5 (record keeping)
// -------------------------------------------------------------------
type jsonRecord struct {
	FileName string            // name of the generated JSON file (or identifier)
	Content  map[string]string // the key/value pairs that were stored
}

// thread‑safe map that holds the relationship between form keys and JSON files
var (
	recordMu sync.RWMutex
	records  = make(map[string]jsonRecord) // map[formKey]jsonRecord
)

// -------------------------------------------------------------------
// 3‑7️⃣  POST handling logic
// -------------------------------------------------------------------
func handlePOST(w http.ResponseWriter, r *http.Request, next http.Handler) {
	const vidiKey = "VIDI"
	const viciKey = "VICI"

	// ----------------------------------------------------------------
	// 3. Look for the VIDI key – if missing, just pass through
	// ----------------------------------------------------------------
	vidiVals, ok := r.Form[vidiKey]
	if !ok || len(vidiVals) == 0 {
		// No VIDI key → nothing special to do, continue down the chain
		next.ServeHTTP(w, r)
		return
	}
	vidiVal := strings.ToUpper(strings.TrimSpace(vidiVals[0]))

	// ----------------------------------------------------------------
	// 3. Validate VIDI value – must be PUSH or PULL (case‑insensitive)
	// ----------------------------------------------------------------
	if vidiVal != "PUSH" && vidiVal != "PULL" {
		http.Error(w, "invalid VIDI value – must be PUSH or PULL", http.StatusBadRequest)
		return
	}

	// ----------------------------------------------------------------
	// 7. Remove any VICI key/value pair from the form before further work
	// ----------------------------------------------------------------
	cleanForm := make(map[string]string)
	for k, vals := range r.Form {
		if strings.EqualFold(k, viciKey) {
			continue // skip VICI entirely
		}
		// Take the first value for simplicity (you can extend to []string if needed)
		if len(vals) > 0 {
			cleanForm[k] = vals[0]
		}
	}

	switch vidiVal {
	case "PUSH":
		// ------------------------------------------------------------
		// 4. Convert the cleaned form data to JSON
		// ------------------------------------------------------------
		jsonBytes, err := json.MarshalIndent(cleanForm, "", "  ")
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		// ------------------------------------------------------------
		// 5. Record which key/value pairs produced which JSON file.
		//    Here we simply store the JSON in memory and give it a
		//    deterministic filename based on a hash of the content.
		// ------------------------------------------------------------
		fileName := generateFileName(jsonBytes) // e.g. "vidi_abc123.json"

		recordMu.Lock()
		records[fileName] = jsonRecord{
			FileName: fileName,
			Content:  cleanForm,
		}
		recordMu.Unlock()

		// ------------------------------------------------------------
		// Respond with the generated JSON (optional – you could also
		// write it to disk if you prefer).
		// ------------------------------------------------------------
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)

	case "PULL":
		// ------------------------------------------------------------
		// 6. Find the stored JSON file that matches the most key/value
		//    pairs with the incoming (cleaned) form data.
		// ------------------------------------------------------------
		bestMatch, matchCount := findBestMatch(cleanForm)

		if bestMatch == nil {
			http.Error(w, "no matching JSON records found", http.StatusNotFound)
			return
		}

		// Return the matched JSON content and a small hint about the match quality
		resp := struct {
			File        string            `json:"file"`
			MatchCount  int               `json:"match_count"`
			TotalStored int               `json:"total_stored"`
			Data        map[string]string `json:"data"`
		}{
			File:        bestMatch.FileName,
			MatchCount:  matchCount,
			TotalStored: len(records),
			Data:        bestMatch.Content,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// -------------------------------------------------------------------
// Utility: deterministic pseudo‑filename from JSON bytes
// -------------------------------------------------------------------
func generateFileName(data []byte) string {
	// Simple checksum – replace with a real hash (sha256) if you need stronger uniqueness
	sum := 0
	for _, b := range data {
		sum = (sum*31 + int(b)) % 0xFFFFFF
	}
	return "vidi_" + strings.ToUpper(strconv.FormatInt(int64(sum), 16)) + ".json"
}

// -------------------------------------------------------------------
// Utility: find the stored JSON record that shares the most key/value pairs
// -------------------------------------------------------------------
func findBestMatch(incoming map[string]string) (*jsonRecord, int) {
	recordMu.RLock()
	defer recordMu.RUnlock()

	var best *jsonRecord
	bestCount := -1

	for _, rec := range records {
		count := countMatchingPairs(incoming, rec.Content)
		if count > bestCount {
			bestCount = count
			tmp := rec // copy to avoid referencing loop variable
			best = &tmp
		}
	}
	if bestCount <= 0 {
		return nil, 0
	}
	return best, bestCount
}

// Count how many exact key/value pairs are identical between two maps
func countMatchingPairs(a, b map[string]string) int {
	matches := 0
	for k, av := range a {
		if bv, ok := b[k]; ok && av == bv {
			matches++
		}
	}
	return matches
}

// -------------------------------------------------------------------
// Optional helper: expose the internal map for external callers
// -------------------------------------------------------------------
func GetRecords() map[string]jsonRecord {
	recordMu.RLock()
	defer recordMu.RUnlock()

	// shallow copy – callers cannot modify the original map
	cpy := make(map[string]jsonRecord, len(records))
	for k, v := range records {
		cpy[k] = v
	}
	return cpy
}

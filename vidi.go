package vidi

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type VidiInterface interface {
	Connect(*http.Request) any
	Commit(any) any
	Create(any) any
	Read(any) any
	Update(any) any
	Delete(any) any
	Find(any) any
	Format(map[string]([]string)) any
	Write(http.ResponseWriter, any) any
}

type VidiContext struct {
	Name        string
	dataManager VidiInterface
}

func InitContext(nameValue string, dm VidiInterface) *VidiContext {
	var vc VidiContext
	vc.Name = nameValue
	if dm == nil {
		dm = InitCore()
	}
	vc.dataManager = dm
	return &vc
}

type VidiCore struct {
	keyMap   map[string]([]uint64)
	valueMap map[string]([]uint64)
	kvMap    map[string]([]uint64)
	fileList []string
}

func (v *VidiCore) Connect(input *http.Request) any {
	return nil

}

func (v *VidiCore) Format(input map[string]([]string)) any {
	return input
}

func (v *VidiCore) Write(w http.ResponseWriter, input any) any {
	if input != nil {
		return ""
	}
	return nil
}

func (v *VidiCore) Commit(input any) any {
	if input != nil {
		return input
	}
	return nil

}

func (v *VidiCore) Create(input any) any {
	value, err := v.parse(input)
	if err != nil {
		return value
	}
	return nil

}

func (v *VidiCore) Read(input any) any {
	value, err := v.parse(input)
	if err != nil {
		return value
	}
	return nil

}

func (v *VidiCore) Update(input any) any {
	value, err := v.parse(input)
	if err != nil {
		return value
	}
	return nil

}

func (v *VidiCore) Delete(input any) any {
	value, err := v.parse(input)
	if err != nil {
		return value
	}
	return nil
}

func (v *VidiCore) Find(input any) any {
	var output []any
	elements, err := v.parse(input)
	if err != nil {
		return elements
	}
	records := make([]uint64, 0)
	for key, val := range elements {
		pairs := make([]string, 0)
		if key == "*" {
			for _, value := range val {
				values, hasValues := v.valueMap[value]
				if hasValues {
					records = append(records, values...)
				}
			}
		} else {
			foundStar := false
			for _, value := range val {
				if value == "*" {
					foundStar = true
					break
				}
				keyval := key + "=" + value
				pairs = append(pairs, keyval)
			}
			if foundStar {
				keys, hasKeys := v.keyMap[key]
				if hasKeys {
					records = append(records, keys...)
				}
			} else {
				for _, pair := range pairs {
					kvs, hasKV := v.kvMap[pair]
					if hasKV {
						records = append(records, kvs...)
					}
				}
			}
		}
	}
	records = v.dedupe(records)
	for i, record := range records {
		fmt.Println(i)
		fid, err := v.identify(record)
		if err != nil {
			fmt.Println(err)
			fmt.Println(i)
			fmt.Println(record)
		}
		readFile := v.Read(fid)
		output = append(output, readFile)
	}
	return output
}

func (v *VidiCore) parse(input any) (map[string]([]string), error) {
	switch v := input.(type) {
	case map[string]([]string):
		return v, nil
	}
	return nil, errors.New("not parsable")
}

func (v *VidiCore) identify(input any) (uint64, error) {
	switch v := input.(type) {
	case uint64:
		return v, nil
	}
	return 0, errors.New("not parsable")
}

func (v *VidiCore) dedupe(slice []uint64) []uint64 {
	seen := make(map[uint64]bool)
	result := []uint64{}

	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

func InitCore() *VidiCore {
	var vc VidiCore
	vc.keyMap = make(map[string]([]uint64))
	vc.valueMap = make(map[string]([]uint64))
	vc.kvMap = make(map[string]([]uint64))
	vc.fileList = make([]string, 0)
	return &vc
}

func (v *VidiContext) Comply(r *http.Request) bool {
	base := path.Base(r.URL.Path)
	switch strings.ToLower(base) {
	case v.Name:
		return true
	default:
		return false
	}
}

func (v *VidiContext) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	v.dataManager.Connect(r)
	formMap := make(map[string]([]string), 0)
	errForm := r.ParseMultipartForm(1024)
	c_Record := false
	r_Record := false
	u_Record := false
	d_Record := false
	var requestValues any
	if errForm != nil {
		fmt.Println(errForm)
	}
	switch r.Method {
	case http.MethodGet:
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, "Failed to parse URL", http.StatusBadRequest)
			return
		}
		// Extract the query parameters into a map
		params := u.Query()
		// Iterate through the parameters and print them
		for key, value := range params {
			//formKeys = append(formKeys, key)
			formMap[key] = value
		}
		r_Record = true
	case http.MethodPost:
		//Find all the elements from the form
		params := r.PostForm
		for key, value := range params {
			//formKeys = append(formKeys, key)
			formMap[key] = value
		}
	default:
		http.Error(w, "Failed to parse method", http.StatusBadRequest)
		return
	}
	requestValues = v.dataManager.Format(formMap)
	//I have a full array of elements
	var output any
	if r_Record {
		output = v.dataManager.Find(requestValues)
	} else {
		if c_Record {
			output = v.dataManager.Create(requestValues)
		} else {
			if d_Record {
				output = v.dataManager.Delete(requestValues)
			} else {
				if u_Record {
					output = v.dataManager.Update(requestValues)
				} else {
					http.Error(w, "Failed to parse info for edit", http.StatusBadRequest)
					return
				}
			}
		}
	}
	output = v.dataManager.Commit(output)
	if output != nil {
		msg := v.dataManager.Write(w, output)
		fmt.Println(msg)
	}
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

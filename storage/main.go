package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tilt-dev/pixeltilt/storage/api"

	"github.com/peterbourgon/diskv/v3"
)

var d = diskv.New(diskv.Options{
	BasePath:     "diskv",
	Transform:    func(s string) []string { return []string{} },
	CacheSizeMax: 1024 * 1024, // 1MB
})

func main() {
	fmt.Println("Storage running!")
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.HandleFunc("/write", write)
	http.HandleFunc("/read", read)
	http.HandleFunc("/list", list)
	http.HandleFunc("/flush", flush)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func read(w http.ResponseWriter, r *http.Request) {
	var rreq api.ReadRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&rreq)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding request: %v", err), http.StatusBadRequest)
		return
	}

	if rreq.Name == "" {
		http.Error(w, "no name specified", http.StatusBadRequest)
		return
	}

	image, err := d.Read(rreq.Name)
	if err != nil {
		// TODO: depending on err, this might be a 404 or something
		http.Error(w, fmt.Sprintf("error reading image from storage: %v", err), http.StatusInternalServerError)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(string(image))
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding image from storage: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")

	rresp := api.ReadResponse{Body: decoded}
	err = json.NewEncoder(w).Encode(rresp)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func write(w http.ResponseWriter, r *http.Request) {
	var wr api.WriteRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&wr)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	if wr.Name == "" {
		http.Error(w, "no Name specified", http.StatusBadRequest)
		return
	}

	if len(wr.Body) == 0 {
		http.Error(w, "no file specified", http.StatusBadRequest)
		return
	}

	filenameWithoutExtension := strings.TrimSuffix(wr.Name, filepath.Ext(wr.Name))
	name := fmt.Sprintf("%s-%s.png", filenameWithoutExtension, time.Now().Format("2006-01-02-15-04-05"))

	err = d.Write(name, wr.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := api.WriteResponse{Name: name}
	err = json.NewEncoder(w).Encode(&response)
}

func list(w http.ResponseWriter, r *http.Request) {
	var lr api.ListResponse
	for key := range d.Keys(nil) {
		lr.Names = append(lr.Names, key)
	}

	err := json.NewEncoder(w).Encode(&lr)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func flush(w http.ResponseWriter, r *http.Request) {
	dir, err := os.ReadDir("/app/diskv")
	if err == nil {
		for _, d := range dir {
			_ = os.RemoveAll(filepath.Join("diskv", d.Name()))
		}
	} else if !os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("error reading diskv directory: %v", err), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("Database entries deleted!\n"))
}

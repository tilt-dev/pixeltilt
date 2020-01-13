package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterbourgon/diskv"
)

var d = diskv.New(diskv.Options{
	BasePath:     "diskv",
	Transform:    func(s string) []string { return []string{} },
	CacheSizeMax: 1024 * 1024, // 1MB
})

func main() {
	fmt.Println("\nStarting up!")
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/access/", access)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	// show index.html
	absPath, err := filepath.Abs("./index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(file)

	// show current db entries
	list := "<br>"
	for key := range d.Keys(nil) {
		list += "<a href='/access/" + key + "/'>" + key + "</a><br>"
	}
	w.Write([]byte(list))

	// TODO: use proper templating? maybe?
}

func access(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.Path, "/")[2]
	image, err := d.Read(key)
	if err != nil {
		fmt.Println(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(image))
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(decoded)

}

func upload(w http.ResponseWriter, r *http.Request) {
	// receive file
	file, handler, err := r.FormFile("myFile")
	defer file.Close()
	if err != nil {
		handleHTTPErr(w, "Error retrieving the file")
		return
	}
	fmt.Printf("Uploaded File: %+v\tFile Size: %+v\tMIME: %+v\n", handler.Filename, handler.Size, handler.Header)

	// create a temp file
	tempFile, err := ioutil.TempFile("", "upload-*.png")
	defer tempFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// save to temp file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		handleHTTPErr(w, fmt.Sprintf("Error reading uploaded file: %v", err))
		return
	}
	tempFile.Write(fileBytes)
	fmt.Println(tempFile.Name())
	defer os.Remove(tempFile.Name())

	// enhance!
	name := tempFile.Name()
	enhanced, err := enhance(name)
	defer os.Remove(enhanced)
	if err != nil {
		handleHTTPErr(w, fmt.Sprintf("Error enhancing %s: %v", name, err))
		return
	}

	// serve output
	output, err := os.Open(enhanced)
	if err != nil {
		handleHTTPErr(w, fmt.Sprintf("Error opening %s: %v", enhanced, err))
		return
	}

	fileBytes, err = ioutil.ReadAll(output)
	w.Header().Set("Content-Type", "image/png")
	w.Write(fileBytes)

	// save to db
	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	d.Write(time.Now().Format("2006-01-02-15-04-05"), []byte(encoded))
	output.Close()

	// TODO: do all of this in memory instead of writing temp files left and right
}

func enhance(file string) (string, error) {
	inputFile := fmt.Sprintf("image=@%s;type=image/png", file)
	outputFile, err := ioutil.TempFile("", "enhanced-*.png")
	if err != nil {
		return "", err
	}
	outputFile.Close()
	// using curl because I forgot how to do this in native go
	cmdOutput, err := exec.Command("curl", "-X", "POST", "http://localhost:5000/model/predict", "-H", "accept: application/json", "-H", "Content-Type: multipart/form-data", "-F", inputFile, "--output", outputFile.Name(), "-s").CombinedOutput()
	if err != nil {
		return "", err
	}
	fmt.Printf("curl output: %s\n", cmdOutput)
	return outputFile.Name(), nil

	// TODO: rewrite in go instead of relying on curl
}

func handleHTTPErr(w http.ResponseWriter, errMsg string) {
	fmt.Println(errMsg)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errMsg))
}

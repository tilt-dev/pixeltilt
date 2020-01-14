package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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
	defer os.Remove(tempFile.Name())
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
	imageType := http.DetectContentType(fileBytes)
	if imageType != "image/png" {
		// https://www.bennadel.com/blog/2434-http-status-codes-for-invalid-data-400-vs-422.htm
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintf("Invalid image type: expected \"image/png\", got: %s", imageType)))
		return
	}
	tempFile.Write(fileBytes)
	fmt.Println(tempFile.Name())

	// enhance!
	name := tempFile.Name()
	enhanced, err := enhance(name)
	if err != nil {
		handleHTTPErr(w, fmt.Sprintf("Error enhancing %s: %v", name, err))
		return
	}

	// serve output
	w.Header().Set("Content-Type", "image/png")
	w.Write(enhanced)

	// save to db
	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	d.Write(time.Now().Format("2006-01-02-15-04-05"), []byte(encoded))
}

func enhance(file string) ([]byte, error) {
	return sendPostRequest("http://localhost:5000/model/predict", file, "image/png")
}

func handleHTTPErr(w http.ResponseWriter, errMsg string) {
	fmt.Println(errMsg)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errMsg))
}

func sendPostRequest(url string, filename string, filetype string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("accept", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

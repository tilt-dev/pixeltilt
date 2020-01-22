package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/windmilleng/enhance/02_storage_service/storage/client"
)

var storage client.Storage

func main() {
	fmt.Println("\nStarting up!")
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	var err error
	storage, err = client.NewStorageClient("http://localhost:8081")
	if err != nil {
		log.Fatalf("initializing storage client: %v", err)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/access/", access)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	// show index.html
	absPath, err := filepath.Abs("./frontend/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	imageKeys, err := storage.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	// show current db entries
	list := "<br>"
	for _, key := range imageKeys {
		list += "<a href='/access/" + key + "/'>" + key + "</a><br>"
	}
	_, err = w.Write([]byte(list))
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: use proper templating? maybe?
}

func access(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.Path, "/")[2]
	image, err := storage.Read(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(image)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	// receive file
	file, header, err := r.FormFile("myFile")
	if err != nil {
		handleHTTPErr(w, "Error retrieving the file")
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\tFile Size: %+v\tMIME: %+v\n", header.Filename, header.Size, header.Header)

	// save to temp file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		handleHTTPErr(w, fmt.Sprintf("Error reading uploaded file: %v", err))
		return
	}
	imageType := http.DetectContentType(fileBytes)
	if imageType != "image/png" {
		// https://www.bennadel.com/blog/2434-http-status-codes-for-invalid-data-400-vs-422.htm
		http.Error(w, fmt.Sprintf("Invalid image type: expected \"image/png\", got: %s", imageType), http.StatusUnprocessableEntity)
		return
	}

	// enhance!
	enhanced, err := enhance(header.Filename, fileBytes)
	if err != nil {
		handleHTTPErr(w, fmt.Sprintf("Error enhancing %s: %v", header.Filename, err))
		return
	}

	// serve output
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(enhanced)
	if err != nil {
		fmt.Println(err)
		return
	}

	// save to db
	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	err = storage.Write(header.Filename, []byte(encoded))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func enhance(name string, image []byte) ([]byte, error) {
	return sendPostRequest("http://localhost:5000/model/predict", name, image)
}

func handleHTTPErr(w http.ResponseWriter, errMsg string) {
	fmt.Println(errMsg)
	http.Error(w, errMsg, http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
}

func sendPostRequest(url string, name string, image []byte) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", name)
	if err != nil {
		return nil, err
	}

	_, err = part.Write(image)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

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

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/windmilleng/enhance/render/api"

	"github.com/windmilleng/enhance/storage/client"
)

var storage client.Storage

func main() {
	fmt.Println("\nStarting up!")
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	var err error
	storage, err = client.NewStorageClient("http://storage")
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

type imageFilter struct {
	url                string
	needsOriginalImage bool
}

var imageFilters = map[string]imageFilter{
	"glitch":     {"http://glitch", false},
	"red":        {"http://red", false},
	"rectangler": {"http://rectangler", true},
}

func upload(w http.ResponseWriter, r *http.Request) {
	originalImageBytes, filename, err := fileFromRequest(r)
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		handleHTTPErr(w, httpStatusError{http.StatusBadRequest, err})
		return
	}

	filters := filtersFromValues(r.PostForm)
	if len(filters) == 0 {
		filters = []string{"glitch", "red", "rectangler"}
	}

	modifiedImage, err := applyFilters(originalImageBytes, filters)
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	// serve output
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(modifiedImage)
	if err != nil {
		fmt.Println(err)
		return
	}

	// save to db
	encoded := base64.StdEncoding.EncodeToString(modifiedImage)
	err = storage.Write(filename, []byte(encoded))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func filtersFromValues(values url.Values) []string {
	var ret []string
	for paramName := range values {
		fmt.Printf("checking if %s is enabling a filter\n", paramName)
		if !strings.HasPrefix(paramName, "filter_") {
			continue
		}

		name := strings.TrimPrefix(paramName, "filter_")
		if _, ok := imageFilters[name]; ok {
			ret = append(ret, name)
		}
	}

	fmt.Printf("returning filter names %v\n", ret)

	return ret
}

type httpStatusError struct {
	code int
	err  error
}

func (e httpStatusError) Error() string {
	return e.err.Error()
}

func fileFromRequest(r *http.Request) (image []byte, filename string, err error) {
	// receive file
	file, header, err := r.FormFile("myFile")
	if err != nil {
		return nil, "", errors.Wrap(err, "getting file from request")
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\tFile Size: %+v\tMIME: %+v\n", header.Filename, header.Size, header.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", errors.Wrap(err, "reading uploaded file from request")
	}

	imageType := http.DetectContentType(fileBytes)
	if imageType != "image/png" {
		// https://www.bennadel.com/blog/2434-http-status-codes-for-invalid-data-400-vs-422.htm
		return nil, "", httpStatusError{http.StatusUnprocessableEntity, fmt.Errorf("invalid image type: expected \"image/png\", got: %s", imageType)}
	}

	return fileBytes, header.Filename, nil
}

func applyFilters(imageBytes []byte, filterNames []string) ([]byte, error) {
	currentImageBytes := append([]byte{}, imageBytes...)

	for _, f := range filterNames {
		var err error
		currentImageBytes, err = applyFilter(imageFilters[f], currentImageBytes, imageBytes)
		if err != nil {
			return nil, fmt.Errorf("Error enhancing %s: %v", f, err)
		}
	}

	return currentImageBytes, nil
}

func applyFilter(filter imageFilter, imageBytes []byte, originalImageBytes []byte) ([]byte, error) {
	rr := api.RenderRequest{Image: imageBytes}
	if filter.needsOriginalImage {
		rr.OriginalImage = originalImageBytes
	}

	resp, err := api.PostRequest(rr, filter.url)
	if err != nil {
		return nil, err
	}

	return resp.Image, nil
}

func handleHTTPErr(w http.ResponseWriter, err error) {
	fmt.Println(err.Error())
	status := http.StatusInternalServerError
	if se, ok := err.(httpStatusError); ok {
		status = se.code
	}
	http.Error(w, err.Error(), status)
}

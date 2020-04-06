package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/windmilleng/enhance/render/api"

	"github.com/windmilleng/enhance/storage/client"
)

type filter struct {
	Label         string `json:"label"`
	URL           string `json:"url"`
	NeedsOriginal bool   `json:"needsoriginal"`
}

// Order matters!
var enabledFilters = []filter{
	filter{"Red", "http://red:8080", false},
	filter{"Glitch", "http://glitch:8080", false},
	filter{"Rectangler", "http://rectangler:8080", true},
}

var storage client.Storage

func main() {
	Benchmark()
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	var err error
	storage, err = client.NewStorageClient("http://storage:8080")
	if err != nil {
		log.Fatalf("initializing storage client: %v", err)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/filters", filters)
	http.HandleFunc("/images", images)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/access/", access)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Document</title>
  </head>
  <body>
    <form
      enctype="multipart/form-data"
      action="/upload"
      method="post"
    >`

	imageKeys, err := storage.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	// <input type="checkbox" id="glitch" name="filter_glitch" checked /><label for="glitch">Glitch</label><br>
	for i := 0; i < len(enabledFilters); i++ {
		lowerName := strings.ToLower(enabledFilters[i].Label)
		html += `<input type="checkbox" id="` + lowerName + `" name="filter_` + lowerName + `" checked /><label for="` + lowerName + `">` + enabledFilters[i].Label + `</label><br>`
	}

	html += `
      <input type="file" name="myFile" />
      <input type="submit" value="upload" />
	</form>
	<br>`

	// show current db entries
	for _, key := range imageKeys {
		html += "<a href='/access/" + key + "/'>" + key + "</a><br><br>"
	}

	html += `
  </body>
</html>`

	_, err = w.Write([]byte(html))
	if err != nil {
		fmt.Println(err)
		return
	}
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

func filters(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(enabledFilters, "", "  ")
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(j)
	if err != nil {
		fmt.Println(err)
	}
}

func images(w http.ResponseWriter, r *http.Request) {
	imageKeys, err := storage.List()
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	j, err := json.MarshalIndent(imageKeys, "", "  ")
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		fmt.Println(err)
	}
}

type uploadResponse struct {
	Name string `json:"name"`
}

// TODO(dmiller): this should return the image URL instead of the image itself
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

	modifiedImage, err := applyFilters(originalImageBytes, filters)
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	// save to db
	encoded := base64.StdEncoding.EncodeToString(modifiedImage)
	name, err := storage.Write(filename, []byte(encoded))
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	resp := uploadResponse{
		Name: name,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error JSON encoding response: %v", err)
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

		for i := 0; i < len(enabledFilters); i++ {
			if strings.ToLower(enabledFilters[i].Label) == name {
				ret = append(ret, enabledFilters[i].Label)
			}
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
		for i := 0; i < len(enabledFilters); i++ {
			if f == enabledFilters[i].Label {
				fmt.Println("APPLYFILTER:", f)
				currentImageBytes, err = applyFilter(enabledFilters[i], currentImageBytes, imageBytes)
				if err != nil {
					return nil, fmt.Errorf("Error enhancing %s: %v", f, err)
				}
			}
		}
	}
	return currentImageBytes, nil
}

func applyFilter(filter filter, imageBytes []byte, originalImageBytes []byte) ([]byte, error) {
	rr := api.RenderRequest{Image: imageBytes}
	if filter.NeedsOriginal {
		rr.OriginalImage = originalImageBytes
	}

	resp, err := api.PostRequest(rr, filter.URL)
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

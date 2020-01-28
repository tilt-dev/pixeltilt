package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type RenderRequest struct {
	Image []byte
}

type RenderReply struct {
	Image []byte
}

type Renderer func(req RenderRequest) (RenderReply, error)

func HttpRenderHandler(renderer Renderer) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rr, err := ReadRequest(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		resp, err := renderer(rr)
		if err != nil {
			msg := fmt.Sprintf("error transforming image: %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		}

		err = WriteResponse(w, resp)
		if err != nil {
			log.Printf("error writing response: %v", err)
		}
	}
}

func ReadRequest(req *http.Request) (RenderRequest, error) {
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	var rr RenderRequest
	err := d.Decode(&rr)
	return rr, errors.Wrap(err, "decoding request body")
}

func WriteResponse(w http.ResponseWriter, resp RenderReply) error {
	return errors.Wrap(json.NewEncoder(w).Encode(resp), "encoding response body")
}

func PostRequest(req RenderRequest, url string) (RenderReply, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return RenderReply{}, errors.Wrap(err, "encoding request body")
	}

	resp, err := http.Post(url, "image/png", &buf)
	if err != nil {
		return RenderReply{}, errors.Wrap(err, "making post request")
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return RenderReply{}, errors.Wrapf(err, "reading response body w/ status %s", resp.Status)
		}
		return RenderReply{}, fmt.Errorf("post request returned status %s: %s", resp.Status, string(body))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RenderReply{}, err
	}

	var reply RenderReply
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	err = d.Decode(&reply)
	return reply, errors.Wrap(err, "decoding reply")
}

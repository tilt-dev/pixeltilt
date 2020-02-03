package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/windmilleng/enhance/storage/api"

	"github.com/pkg/errors"
)

type Storage struct {
	url *url.URL
}

func NewStorageClient(rawurl string) (Storage, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return Storage{}, errors.Wrapf(err, "parsing url %q", rawurl)
	}
	return Storage{url: u}, nil
}

func (sc *Storage) Write(name string, b []byte) (string, error) {
	wr := api.WriteRequest{Name: name, Body: b}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(wr)
	if err != nil {
		return "", errors.Wrap(err, "encoding write request")
	}
	u, _ := url.Parse(sc.url.String())
	u.Path = "/write"
	resp, err := http.Post(u.String(), "application/json", buf)
	if err != nil {
		return "", errors.Wrap(err, "http posting")
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Wrap(err, "reading response body on non-200 response")
		}
		return "", errors.Errorf("http status %s, body %q", resp.Status, string(respBody))
	}

	var wresp api.WriteResponse
	err = json.NewDecoder(resp.Body).Decode(&wresp)
	if err != nil {
		errors.Wrap(err, "reading response body on 200 response")
	}

	return wresp.Name, nil
}

func (sc *Storage) Read(name string) ([]byte, error) {
	rreq := api.ReadRequest{Name: name}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(rreq)
	if err != nil {
		return nil, errors.Wrap(err, "encoding read request")
	}

	u, _ := url.Parse(sc.url.String())
	u.Path = "/read"
	resp, err := http.Post(u.String(), "application/json", buf)
	if err != nil {
		return nil, errors.Wrap(err, "http posting")
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "reading response body on non-200 response")
		}
		return nil, errors.Errorf("http status %s, body %q", resp.Status, string(respBody))
	}

	var rresp api.ReadResponse
	err = json.NewDecoder(resp.Body).Decode(&rresp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding http response")
	}

	return rresp.Body, nil
}

func (sc *Storage) List() ([]string, error) {
	u, _ := url.Parse(sc.url.String())
	u.Path = "/list"
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "http getting")
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "reading response body on non-200 response")
		}
		return nil, errors.Errorf("http status %s, body %q", resp.Status, string(respBody))
	}

	var lr api.ListResponse
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&lr)
	if err != nil {
		return nil, errors.Wrap(err, "decoding http response")
	}

	// Don't return null because null is the root of all evil
	if lr.Names == nil {
		return []string{}, nil
	}

	return lr.Names, nil
}

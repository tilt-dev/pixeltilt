package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestRed(t *testing.T) {
	input, err := os.Open("../frontend/public/duck.png")
	if err != nil {
		panic(err)
	}
	pic, err := png.Decode(input)
	processed := red(pic)
	var got bytes.Buffer
	err = png.Encode(&got, processed)
	if err != nil {
		panic(err)
	}
	wanted, err := ioutil.ReadFile("test/red.png")
	if err != nil {
		panic(err)
	}
	if bytes.Compare(got.Bytes(), wanted) != 0 {
		t.Errorf("red filter got != wanted")
	}
}

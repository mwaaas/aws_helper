package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type JSON interface{}

func getRequestBody(request *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		return nil, err
	}
	// Because go lang is a pain in the ass if you read the body then any susequent calls
	// are unable to read the body again....
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, err
}

// Get a json decoder for a given requests body
func requestBodyDecoder(request *http.Request) (*json.Decoder, error) {
	// Read body to buffer
	body, err := getRequestBody(request)
	if err != nil {
		return nil, err
	}

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body))), nil
}

// Parse the requests body
func parseRequestBodyAsMap(request *http.Request) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	decoder, err := requestBodyDecoder(request)

	if err != nil {
		return nil, err
	}
	err = decoder.Decode(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

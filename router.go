package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// New returns a new router
func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/sns/send_message", sqsSendMessage)
	return r
}

func sqsSendMessage(res http.ResponseWriter, req *http.Request) {
	method := req.URL.Query().Get("method")
	reqBody, _ := parseRequestBodyAsMap(req)
	contextLogger := log.WithFields(
		log.Fields{
			"url":          req.URL,
			"method":       method,
			"request_body": reqBody,
		})

	contextLogger.Info("Handling URL request")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	endpoint := &sess.Config.Endpoint

	if reqEndPoint, ok := reqBody["endpoint"]; ok {
		if reqEndPointStr, ok := reqEndPoint.(string); ok {
			if reqEndPoint != "" {
				*endpoint = aws.String(reqEndPointStr)
			}
		} else {
			panic("endpoint should be string")
		}

	}

	svc := sns.New(sess)

	// if message structure is in json we convert message to json
	if value, ok := reqBody["MessageStructure"]; ok {
		message := reqBody["message"]
		if messageMap, ok := message.(map[string]interface{}); ok {
			if value == "json" {
				if _, ok := messageMap["default"]; !ok {
					messageMap["default"] = "default was not set"
				}
			}
			messageBytes, _ := json.Marshal(messageMap)
			reqBody["message"] = aws.String(string(messageBytes))
		}
	}
	publishInput := &sns.PublishInput{}
	err := mapstructure.Decode(reqBody, publishInput)

	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
		_, err := res.Write([]byte(err.Error()))
		if err != nil {
			panic(err)
		}
		return
	}

	resp, err := svc.Publish(publishInput)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
		_, err := res.Write([]byte(err.Error()))

		if err != nil {
			panic(err)
		}
	} else {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(resp.GoString()))

		if err != nil {
			panic(err)
		}
	}
	return
}

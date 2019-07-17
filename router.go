package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// New returns a new router
func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/sns/send_message", snsSendMessage)
	r.HandleFunc("/sqs/send_message", sqsSendMessage)
	return r
}


func getAwsSession(endpoint string)  *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))



	awsEndpoint := &sess.Config.Endpoint

	if endpoint != "" {
		*awsEndpoint = aws.String(endpoint)
	}


	return sess
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
	svc := sqs.New(getAwsSession(reqBody["endpoint"].(string)))

	messageBody := reqBody["message"]
	messageBytes, err := json.Marshal(messageBody)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
		_, err := res.Write([]byte(err.Error()))
		if err != nil {
			panic(err)
		}
		return
	}

	reqBody["MessageBody"] = aws.String(string(messageBytes))

	messageParam := &sqs.SendMessageInput{}

	err = mapstructure.Decode(reqBody, messageParam)

	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
		_, err := res.Write([]byte(err.Error()))
		if err != nil {
			panic(err)
		}
		return
	}

	resp, err := svc.SendMessage(messageParam)

	if err != nil {
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

func snsSendMessage(res http.ResponseWriter, req *http.Request) {
	method := req.URL.Query().Get("method")
	reqBody, _ := parseRequestBodyAsMap(req)
	contextLogger := log.WithFields(
		log.Fields{
			"url":          req.URL,
			"method":       method,
			"request_body": reqBody,
		})

	contextLogger.Info("Handling URL request")


	svc := sns.New(getAwsSession(reqBody["endpoint"].(string)))

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

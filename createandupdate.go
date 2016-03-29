package main

import (
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"math/rand"
)

var logger *log.Logger
var client *http.Client

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func login(url string, appID string, appKey string, clientID string, clientSecret string) string {

	type oauthRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	type oauthResponse struct {
		AccessToken string `json:"access_token"`
	}

	requestBody := oauthRequest{ClientID: clientID, ClientSecret: clientSecret, }
	jsonRequest, _ := json.Marshal(requestBody)

	logger.Printf("Log-in request is %s", string(jsonRequest))

	req, err := http.NewRequest("POST", url + "/oauth2/token", bytes.NewReader(jsonRequest))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Add("x-kii-appid", appID)
	req.Header.Add("x-kii-appkey", appKey)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Fatal(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if (res.StatusCode != 200) {
		logger.Fatal("Log-in failed ", res.Status, string(body))
	}
	var response oauthResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Printf("Logged-in (token=%s)", response.AccessToken)
	return response.AccessToken

}

func updateObject(url string, token string, appID string, appKey string, objectID string) {

	type objectUpdateRequest struct {
		FirstName      string
		SomeOtherStuff string
	}

	type objectUpdateResponse struct {
		ModifiedAt  int64 `json:"modifiedAt"`
		CreatedAt   int64 `json:"createdAt"`
		EntityTagID string `json:"entityTagID"`
	}

	client := &http.Client{}
	jsonRequest, _ := json.Marshal(objectUpdateRequest{
		FirstName: randomString(20),
		SomeOtherStuff: randomString(45),
	})
	req, err := http.NewRequest("PATCH", url + "/apps/" + appID + "/buckets/testing_kii/objects/" + objectID, bytes.NewReader(jsonRequest))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "bearer " + token)

	res, err := client.Do(req)
	if err != nil {
		logger.Fatal(err.Error())
	}
	if (res.StatusCode != 200) {
		logger.Fatal("Update object failed ", res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Fatal(err.Error())
	}
	var response objectUpdateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Printf("Updated object (response=%s)", string(body))
}

func createObject(url string, token string, appID string, appKey string) string {

	type objectCreationRequest struct {
		FirstName    string
		LastName     string
		EmailAddress string
	}

	type objectCreationResponse struct {
		ObjectID    string `json:"objectID"`
		CreatedAt   int64 `json:"createdAt"`
		EntityTagID string `json:"entityTagID"`
		DataType    string `json:"dataType"`
	}
	client := &http.Client{}
	jsonRequest, _ := json.Marshal(objectCreationRequest{
		FirstName: randomString(20),
		LastName: randomString(30),
		EmailAddress: randomString(25),
	})

	req, err := http.NewRequest("POST", url + "/apps/" + appID + "/buckets/testing_kii/objects", bytes.NewReader(jsonRequest))
	if err != nil {
		logger.Fatal(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "bearer " + token)

	res, err := client.Do(req)
	if err != nil {
		logger.Fatal(err.Error())
	}
	if (res.StatusCode != 201) {
		logger.Fatal("Create object failed ", res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	var response objectCreationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Printf("Created object (id=%s)", response.ObjectID)
	return response.ObjectID
}

func init() {
	logger = log.New(os.Stdout, "createandupdate - ", log.LstdFlags | log.Lmicroseconds)
	client = &http.Client{}
}

func main() {
	if len(os.Args) != 6 {
		log.Fatal("Incorrect usage. Usage: url appID appKey clientID clientSecret")
	}
	url := os.Args[1]
	appID := os.Args[2]
	appKey := os.Args[3]
	clientID := os.Args[4]
	clientSecret := os.Args[5]
	logger.Printf("Settings: url=%s, appID=%s, appKey=%s, clientID=%s, clientSecret=%s", url, appID, appKey, clientID, clientSecret)
	token := login(url, appID, appKey, clientID, clientSecret)
	for {
		updateObject(url, token, appID, appKey, createObject(url, token, appID, appKey))
	}
}
package main

import (
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"math/rand"
	"time"
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

func register(url string, appID string, appKey string) (userID string, token string) {
	type regAndAuthRequest struct {
		UserName string `json:"loginName"`
		Password string `json:"password"`
	}

	type regAndAuthResponse struct {
		UserID      string `json:"userID"`
		AccessToken string `json:"_accessToken"`
	}

	requestBody := regAndAuthRequest{UserName: randomString(10), Password: randomString(10), }
	jsonRequest, _ := json.Marshal(requestBody)

	logger.Printf("Registration request is %s", string(jsonRequest))

	req, err := http.NewRequest("POST", url + "/apps/"+appID+"/users", bytes.NewReader(jsonRequest))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Add("x-kii-appid", appID)
	req.Header.Add("x-kii-appkey", appKey)
	req.Header.Add("content-type", "application/vnd.kii.RegistrationAndAuthorizationRequest+json")
	req.Header.Add("accept", "application/vnd.kii.RegistrationAndAuthorizationResponse+json")

	res, err := client.Do(req)
	if err != nil {
		logger.Fatal(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if (res.StatusCode != 201) {
		logger.Fatal("Registration failed ", res.Status, string(body))
	}
	var response regAndAuthResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Printf("Registered (response=%s)", string(body))
	return response.UserID, response.AccessToken
}

func userLogin(url string, appID string, appKey string, userName string, password string) string {
	type oauthRequest struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	type oauthResponse struct {
		AccessToken string `json:"access_token"`
	}

	requestBody := oauthRequest{UserName: userName, Password: password, }
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
	logger.Printf("Logged-in as admin (token=%s)", response.AccessToken)
	return response.AccessToken
}
func adminLogin(url string, appID string, appKey string, clientID string, clientSecret string) string {

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
	logger.Printf("Logged-in as admin (token=%s)", response.AccessToken)
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
	req, err := http.NewRequest("PATCH", url + "/apps/" + appID + "/users/me/buckets/testing_kii/objects/" + objectID, bytes.NewReader(jsonRequest))
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

	req, err := http.NewRequest("POST", url + "/apps/" + appID + "/users/me/buckets/testing_kii/objects", bytes.NewReader(jsonRequest))
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
	rand.Seed(time.Now().UTC().UnixNano())
	logger = log.New(os.Stdout, "createandupdate - ", log.LstdFlags | log.Lmicroseconds)
	client = &http.Client{}
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Incorrect usage. Usage: url appID appKey")
	}
	url := os.Args[1]
	appID := os.Args[2]
	appKey := os.Args[3]
	logger.Printf("Settings: url=%s, appID=%s, appKey=%s", url, appID, appKey)
	_, token := register(url, appID, appKey)
	for {
		updateObject(url, token, appID, appKey, createObject(url, token, appID, appKey))
	}
}
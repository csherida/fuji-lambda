package app

import (
	"encoding/json"
	"fuji-alexa/internal/models/fuji"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var appleUserToken string

func getAppleUserToken(amazonToken string) string {

	if appleUserToken != "" {
		return appleUserToken
	}

	url := "https://ff7lyzbjr9.execute-api.us-east-1.amazonaws.com/prod/fujiaccount?amazon-token=" + amazonToken

	// Create a Bearer string by appending string access token
	var apiKey = getSecret("FujiAccountAPIKey")
	if apiKey == "" {
		log.Println("Fuji Account API Key is blank.")
		//TODO: Add error return
		return ""
	}

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("x-api-key", apiKey)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	//TODO: Handle 401 errors
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	log.Println("Length of body response: " + strconv.Itoa(len(body)))

	var responseObject models.FujiAccount
	json.Unmarshal(body, &responseObject)

	// Cache value
	appleUserToken = responseObject.AppleToken

	//TODO: Handle nulls and empty strings
	return responseObject.AppleToken
}

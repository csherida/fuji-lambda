package app

import (
	"encoding/json"
	"errors"
	"fuji-alexa/internal/models/fuji"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetFujiAccount(amazonToken string) (*models.FujiAccount, error) {

	url := "https://ff7lyzbjr9.execute-api.us-east-1.amazonaws.com/prod/fujiaccount?amazon-token=" + amazonToken

	// Create a Bearer string by appending string access token
	var apiKey = getSecret("FujiAccountAPIKey")
	if apiKey == "" {
		log.Printf("Unable to get Fuji API Key from secrets management.")
		return nil, errors.New("failed to get Fuji API Key")
	}

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("x-api-key", apiKey)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return nil, err
	}
	if resp.StatusCode == 404 {
		log.Println("Unable to find Fuji account for Amazon user: ", amazonToken)
		return nil, errors.New("user not found")
	} else if resp.StatusCode >= 400 {
		log.Println("Received error code: ", strconv.Itoa(resp.StatusCode))
		return nil, errors.New("error fetching Fuji account")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var responseObject models.FujiAccount
	json.Unmarshal(body, &responseObject)

	//TODO: Handle nulls and empty strings
	return &responseObject, nil
}

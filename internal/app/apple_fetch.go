package app

import (
	"encoding/json"
	"errors"
	"fuji-alexa/internal/models/apple"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func fetchAppleMusicData(amazonToken string, url string) (*apple.AppleResponse, error) {

	// Create a Bearer string by appending string access token
	var secret = getSecret("FujiAppleMusicToken")
	if secret == "" {
		log.Println("Apple Music token is blank.")
		return nil, errors.New("Apple Music token is blank.")
	}
	var bearer = "Bearer " + secret

	// TODO: cache this value
	// Get the Apple User Token associated with this amazon user token
	var appleUserToken = getAppleUserToken(amazonToken)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header and user token to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Music-User-Token", appleUserToken)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	//TODO: Handle 401 errors
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
		return nil, err
	}

	log.Println("Length of body response: " + strconv.Itoa(len(body)))

	var responseObject apple.AppleResponse
	json.Unmarshal(body, &responseObject)
	return &responseObject, nil
}

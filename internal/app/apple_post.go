package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fuji-alexa/internal/models/apple"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func postAppleMusicData(amazonToken string, url string, tracks apple.AppleResponse) error {

	// Create a Bearer string by appending string access token
	var secret = getSecret("FujiAppleMusicToken")
	if secret == "" {
		log.Println("Apple Music token is blank.")
		return errors.New("apple music token is blank")
	}
	var bearer = "Bearer " + secret

	// Get the Apple User Token associated with this amazon user token
	var appleUserToken = getAppleUserToken(amazonToken)

	// Setup request body
	reqBody, err := json.Marshal(tracks)
	if err != nil {
		log.Fatalf("Unable to marshal tracks for Amazon Token %v with URL: %v", amazonToken, url)
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	// add authorization header and user token to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Music-User-Token", appleUserToken)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	//TODO: Handle 401, 404 errors
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
		return err
	}

	log.Println("Length of body response: " + strconv.Itoa(len(body)))

	return nil
}

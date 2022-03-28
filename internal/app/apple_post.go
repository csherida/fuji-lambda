package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fuji-alexa/internal/models/apple"
	"log"
	"net/http"
	"strconv"
)

func postAppleMusicData(amazonToken string, url string, tracks apple.AppleTrackRequest) error {

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

	if (err != nil) || (resp.StatusCode >= 400) {
		log.Println("Error on response.\n[ERROR] -", err)
		if err == nil {
			err = errors.New("Received a status code of " + strconv.Itoa(resp.StatusCode) + " for " + url)
		}
		return err
	}
	defer resp.Body.Close()

	// Assume no response body needs to be parsed for the posting
	return nil
}

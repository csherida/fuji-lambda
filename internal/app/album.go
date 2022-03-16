package app

import (
	"encoding/json"
	"fuji-alexa/internal/models/apple"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetAlbum(id int) Album {

	url := "https://api.music.apple.com/v1/catalog/us/albums/" + strconv.Itoa(id)

	// Create a Bearer string by appending string access token
	var secret = getSecret("FujiAppleMusicToken")
	if secret == "" {
		log.Println("Apple Music token is blank.")
		var album Album
		return album
	}
	var bearer = "Bearer " + secret

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

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

	var responseObject apple.AppleAlbum
	json.Unmarshal(body, &responseObject)

	//TODO: Handle nulls
	var album Album
	album.Name = responseObject.Data[0].Attributes.Name
	album.ArtistName = responseObject.Data[0].Attributes.ArtistName

	return album
}

type Album struct {
	ArtistName string
	Name       string
}

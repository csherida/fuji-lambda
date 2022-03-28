package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fuji-alexa/internal/models/apple"
	models "fuji-alexa/internal/models/fuji"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func AddTracksToPlaylist(amazonToken string, playlistID string, tracks apple.AppleTrackRequest) error {

	url := "https://api.music.apple.com/v1/me/library/playlists/" + playlistID + "/tracks"
	err := postAppleMusicData(amazonToken, url, tracks)
	if err != nil {
		log.Fatalf("Unable to add tracks for Amazon Token %v and playlist %v.", amazonToken, playlistID)
		return err
	}
	return nil
}

func GetPlaylistCount(amazonToken string) int {

	url := "https://api.music.apple.com/v1/me/library/playlists"

	responseObject, err := fetchAppleMusicData(amazonToken, url)
	if err != nil {
		log.Fatalf("Unalbe to get user's playlist for %v\n", err)
		return 0
	}

	return responseObject.Meta.Total
}

func GetPlaylistTracks(amazonToken string, playlistID string, pageOffset ...int) (*apple.AppleResponse, error) {

	// See if pagination is required
	offset := 0
	if len(pageOffset) > 0 {
		offset = pageOffset[0]
	}

	// Construct URL and add an offset for pagination in case needed
	url := "https://api.music.apple.com/v1/me/library/playlists/" + playlistID + "/tracks"
	if offset > 0 {
		url += "?offset=" + strconv.Itoa(offset)
	}

	responseObject, err := fetchAppleMusicData(amazonToken, url)
	if err != nil {
		log.Fatalf("Unable to shuffle playlist for playlist ID: %v", playlistID)
		return nil, err
	}

	trackCount := responseObject.Meta.Total
	log.Println("Number of tracks in the playlist: " + strconv.Itoa(trackCount))

	return responseObject, nil
}

func CreatePlaylist(amazonToken string, origPlaylistName string) (*models.FujiPlaylist, error) {
	// Create a Bearer string by appending string access token
	var secret = getSecret("FujiAppleMusicToken")
	if secret == "" {
		log.Println("Apple Music token is blank.")
		return nil, errors.New("apple music token is blank")
	}
	var bearer = "Bearer " + secret

	// Get the Apple User Token associated with this amazon user token
	var appleUserToken = getAppleUserToken(amazonToken)

	// Setup request body
	rand.Seed(time.Now().UnixNano())
	newPlaylistName := origPlaylistName + " " + strconv.Itoa(rand.Intn(999-1))
	playlistRequest := apple.PlaylistRequest{}
	playlistRequest.Attributes.Name = newPlaylistName
	playlistRequest.Attributes.Description = "Created by Fuji to randomize playlist " + origPlaylistName
	//TODO: Get Folder ID of the Fuji folder (if exists)
	playlistRequest.Relationships.Parent.Data = append(playlistRequest.Relationships.Parent.Data,
		apple.ParentFolderData{ID: "p.KoZ1ACaN0ZO2", Type: "library-playlist-folders"})

	reqBody, err := json.Marshal(playlistRequest)
	if err != nil {
		log.Fatalf("Unable to create a new playlist for %v", newPlaylistName)
	}

	// Create a new request using http
	url := "https://api.music.apple.com/v1/me/library/playlists/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	// add authorization header and user token to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Music-User-Token", appleUserToken)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	//TODO: Handle 401, 404 errors
	if (err != nil) || (resp.StatusCode >= 400) {
		log.Println("Error on response trying to create a playlist.\n[ERROR] -", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes for playlist creation:", err)
		return nil, err
	}

	log.Println("Length of body response: " + strconv.Itoa(len(body)))
	var responseObject apple.AppleResponse
	json.Unmarshal(body, &responseObject)

	return &models.FujiPlaylist{PlaylistID: responseObject.Data[0].ID, Name: newPlaylistName}, nil
}

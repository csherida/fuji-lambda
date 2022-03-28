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
	"sync"
	"time"
)

// This function takes in the user's utterance in Alexa and looks in the list of playlists for it
// TODO: Perform fuzzy matching
func FindPlaylist(amazonToken string, playlistNameRqst string) string {
	playlistID := ""
	playlists, err := GetPlaylists(amazonToken)
	if err != nil {
		log.Fatalf("Error finding playlist %v", playlistNameRqst)
		return ""
	}
	for _, playlist := range playlists.FujiPlaylist {
		if playlist.Name == playlistNameRqst {
			playlistID = playlist.ID
		}
	}
	return playlistID
}

func GetPlaylistCount(amazonToken string) int {

	url := "https://api.music.apple.com/v1/me/library/playlists"

	responseObject, err := fetchAppleMusicData(amazonToken, url)
	if err != nil {
		log.Fatalf("Unalbe to get user's playlist count for %v\n", err)
		return 0
	}

	return responseObject.Meta.Total
}

func GetPlaylists(amazonToken string) (*models.FujiPlaylists, error) {

	// Make the initial call
	fujiPlaylists, err := getPlaylistsFromApple(amazonToken, 0)
	if err != nil {
		log.Fatalf("Unalbe to get user's playlist for %v\n", err)
		return nil, err
	}

	// Apple Music APIs paginate at 25 playlists
	if fujiPlaylists.OverallPlaylistCount <= 25 {
		return fujiPlaylists, nil
	}

	// If we have more than 25, the let's fan-out to capture them all
	offsetCount := calculateOffset(fujiPlaylists.OverallPlaylistCount, 25)
	var tracksMap sync.Map
	wg := sync.WaitGroup{}

	//Remember, we've already called it once, so this is for subsequent calls
	for i := 1; i < offsetCount; i++ {
		wg.Add(1)
		go func(idx int) {
			newPlaylists, err := getPlaylistsFromApple(amazonToken, idx*25)
			if err != nil {
				log.Fatalf("Unable to rtrieve playlist for offset %v", i)
				panic(err)
			}
			// While slice appending is likely thread-safe, giving a dedicated memory space just in case
			tracksMap.Store(idx, newPlaylists.FujiPlaylist)
			log.Printf("Received playlists for offset %v.", idx)
			wg.Done()
		}(i)
	}
	wg.Wait()

	offsetMapCount := lenSyncMap(&tracksMap)
	log.Printf("Number of offsets captured: %v", offsetMapCount)

	tracksMap.Range(func(key, value interface{}) bool {
		playlistFromMap := value.([]models.FujiPlaylist)
		fujiPlaylists.FujiPlaylist = append(fujiPlaylists.FujiPlaylist, playlistFromMap...)
		return true
	})

	return fujiPlaylists, nil
}

func getPlaylistsFromApple(amazonToken string, offset int) (*models.FujiPlaylists, error) {

	// Construct the URL
	url := "https://api.music.apple.com/v1/me/library/playlists"
	if offset > 0 {
		url = url + "?offset=" + strconv.Itoa(offset)
	}

	// Call Apple's APIs
	responseObject, err := fetchAppleMusicData(amazonToken, url)
	if err != nil {
		log.Fatalf("Unalbe to get user's playlist for %v\n", err)
		return nil, err
	}

	// Scrub off unnecessary data
	var playlists []models.FujiPlaylist
	for _, playlist := range responseObject.Data {
		newPlaylist := models.FujiPlaylist{ID: playlist.ID, Name: playlist.Attributes.Name}
		playlists = append(playlists, newPlaylist)
	}

	// Box up the response and send
	fujiPlaylists := models.FujiPlaylists{FujiPlaylist: playlists, OverallPlaylistCount: responseObject.Meta.Total}
	return &fujiPlaylists, nil
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
	// TODO: Do a literal instantiation of playlist request
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

	if (err != nil) || (resp.StatusCode >= 400) {
		log.Println("Error on response trying to create a playlist.\n[ERROR] -", err)
		if err == nil {
			err = errors.New("Received a status code of " + strconv.Itoa(resp.StatusCode) + " for " + url)
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes for playlist creation:", err)
		return nil, err
	}

	var responseObject apple.AppleResponse
	json.Unmarshal(body, &responseObject)

	return &models.FujiPlaylist{ID: responseObject.Data[0].ID, Name: newPlaylistName}, nil
}

package app

import (
	"fuji-alexa/internal/models/apple"
	"log"
	"strconv"
)

func AddTracksToPlaylist(amazonToken string, playlistID string, tracks apple.AppleResponse) error {

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

func GetPlaylist(amazonToken string, playlistID string) (*apple.AppleResponse, error) {

	url := "https://api.music.apple.com/v1/me/library/playlists/" + playlistID + "/tracks"

	responseObject, err := fetchAppleMusicData(amazonToken, url)
	if err != nil {
		log.Fatalf("Unable to shuffle playlist for playlist ID: %v", playlistID)
		return nil, err
	}

	trackCount := responseObject.Meta.Total
	log.Println("Number of tracks in the playlist: " + strconv.Itoa(trackCount))

	return responseObject, nil
}

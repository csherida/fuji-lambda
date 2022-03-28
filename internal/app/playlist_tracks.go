package app

import (
	"fuji-alexa/internal/models/apple"
	"log"
	"strconv"
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
		log.Printf("Unable to shuffle playlist for playlist ID: %v", playlistID)
		return nil, err
	}

	trackCount := responseObject.Meta.Total
	log.Println("Number of tracks in the playlist: " + strconv.Itoa(trackCount))

	return responseObject, nil
}

package app

import (
	"fuji-alexa/internal/models/apple"
	models "fuji-alexa/internal/models/fuji"
	"log"
	"strconv"
)

func AddTracksToPlaylist(acct *models.FujiAccount, playlistID string, tracks apple.AppleTrackRequest) error {

	url := "https://api.music.apple.com/v1/me/library/playlists/" + playlistID + "/tracks"
	err := postAppleMusicData(acct, url, tracks)
	if err != nil {
		log.Fatalf("Unable to add tracks for Amazon Token %v and playlist %v.", acct.AmazonToken, playlistID)
		return err
	}
	return nil
}

func GetPlaylistTracks(acct *models.FujiAccount, playlistID string, pageOffset ...int) (*apple.AppleResponse, error) {

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

	responseObject, err := fetchAppleMusicData(acct, url)
	if err != nil {
		log.Printf("Unable to shuffle playlist for playlist ID: %v", playlistID)
		return nil, err
	}

	trackCount := responseObject.Meta.Total
	log.Println("Number of tracks in the playlist: " + strconv.Itoa(trackCount))

	return responseObject, nil
}

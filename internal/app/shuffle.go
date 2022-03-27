package app

import (
	"fuji-alexa/internal/models/apple"
	"log"
	"math/rand"
	"time"
)

func ShufflePlaylist(amazonToken string, playlistName string) (string, error) {
	//TODO: Lookup playlist ID by name
	log.Printf("Shuffling playlist %v", playlistName)
	return shuffle(amazonToken, "p.5x1WhOxAz9v")
}

func shuffle(amazonToken string, origPlaylistID string) (string, error) {

	tracks, err := getTracks(amazonToken, origPlaylistID)

	log.Printf("Ordered first track: %v", tracks.Data[0].Attributes.Name)
	log.Printf("Ordered last track: %v", tracks.Data[len(tracks.Data)-1].Attributes.Name)

	// Shuffle the list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tracks.Data), func(i, j int) { tracks.Data[i], tracks.Data[j] = tracks.Data[j], tracks.Data[i] })

	log.Printf("Shuffled first track: %v", tracks.Data[0].Attributes.Name)
	log.Printf("Shuffled last track: %v", tracks.Data[len(tracks.Data)-1].Attributes.Name)

	// TODO: create new playlist, checking that the Fuji folder exists
	newPlaylistID := "p.zpGExIm2pvM5"
	err = AddTracksToPlaylist(amazonToken, newPlaylistID, *tracks)
	if err != nil {
		log.Fatalf("Unable to create new, suffled playlist for: %v", newPlaylistID)
		return "", err
	}

	// TODO: return name of new playlist, will likely have to iterate through all playlists
	return "", nil
}

// This function will get the playlist tracks and scrub it so only IDs are returned
func getTracks(amazonToken string, origPlaylistID string, pageOffset ...int) (*apple.AppleResponse, error) {

	// See if pagination is required
	offset := 0
	if len(pageOffset) > 0 {
		offset = pageOffset[0]
	}

	// Call function to retrieve full data set
	tracks, err := GetPlaylistTracks(amazonToken, origPlaylistID, offset)
	if err != nil {
		log.Fatalf("Unable to shuffle and scrub playlist for playlist ID: %v", origPlaylistID)
		return nil, err
	}

	// Create new Apple object to hold just IDs
	var scrubbedTracks *apple.AppleResponse
	scrubbedTracks = new(apple.AppleResponse)

	// Loop through returned object and pull out just the ID
	for _, track := range tracks.Data {
		var scrubbedTrack apple.Data
		scrubbedTrack.ID = track.ID
		scrubbedTracks.Data = append(scrubbedTracks.Data, scrubbedTrack)
	}

	return scrubbedTracks, nil
}

// This function calculates the pagination and how many times we have to call Apple Music
func calculateOffset(trackCount int) int {
	offsetCount := trackCount / 100
	if trackCount%100 > 0 {
		offsetCount++
	}
	return offsetCount
}

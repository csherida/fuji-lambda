package app

import (
	"fuji-alexa/internal/models/apple"
	"log"
	"math/rand"
	"sync"
	"time"
)

// This function is the primary orchestrator for getting tracks, shuffling and writing to a new Playlist
func ShufflePlaylist(amazonToken string, origPlaylistID string, origPlaylistName string) (string, error) {

	log.Printf("Shuffling playlist %v", origPlaylistID)

	// Get tracks, scrub and shuffle
	tracks, err := getTracks(amazonToken, origPlaylistID)
	scrubbedTracks := scrubTracks(tracks)
	scrubbedTracks = shuffle(scrubbedTracks)

	//Check for pagination.  For Apple, it's playlists with over 100 tracks
	offsetCount := calculateOffset(tracks.Meta.Total, 100)
	if offsetCount > 1 {
		var tracksMap sync.Map
		wg := sync.WaitGroup{}

		//Remember, we've already called it once, so this is for subsequent calls
		for i := 1; i < offsetCount; i++ {
			wg.Add(1)
			go func(idx int) {
				newTracks, err := getTracks(amazonToken, origPlaylistID, idx*100)
				if err != nil {
					log.Fatalf("Unable to rtrieve tracks for offset %v in playlist %v", i, origPlaylistID)
					panic(err)
				}
				scrubbedNewTracks := scrubTracks(newTracks)
				//scrubbedTracks.Data = append(scrubbedTracks.Data, scrubbedNewTracks.Data...)
				tracksMap.Store(idx, scrubbedNewTracks.Data)
				log.Printf("Received tracks for offset %v in playlist %v", i, origPlaylistID)
				wg.Done()
			}(i)
		}
		wg.Wait()

		offsetMapCount := lenSyncMap(&tracksMap)
		log.Printf("Number of offsets captured: %v", offsetMapCount)

		tracksMap.Range(func(key, value interface{}) bool {
			tracks2 := value.([]apple.TrackData)
			scrubbedTracks.Data = append(scrubbedTracks.Data, tracks2...)
			return true
		})
	}

	scrubbedTracks = shuffle(scrubbedTracks)

	// TODO: check that the Fuji folder exists
	// TODO: feed in name of requested playlist
	newPlaylist, err := CreatePlaylist(amazonToken, origPlaylistName)
	if err != nil {
		log.Fatalf("Unable to create a new playlist")
		return "", err
	}
	err = AddTracksToPlaylist(amazonToken, newPlaylist.ID, *scrubbedTracks)

	if err != nil {
		log.Fatalf("Unable to add tracks to new, suffled playlist: %v", newPlaylist.ID)
		return "", err
	}

	return newPlaylist.Name, nil
}

func shuffle(tracks *apple.AppleTrackRequest) *apple.AppleTrackRequest {

	log.Printf("Ordered first track: %v", tracks.Data[0].ID)
	log.Printf("Ordered last track: %v", tracks.Data[len(tracks.Data)-1].ID)

	// Shuffle the list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tracks.Data), func(i, j int) { tracks.Data[i], tracks.Data[j] = tracks.Data[j], tracks.Data[i] })

	log.Printf("Shuffled first track: %v", tracks.Data[0].ID)
	log.Printf("Shuffled last track: %v", tracks.Data[len(tracks.Data)-1].ID)

	// TODO: return name of new playlist, will likely have to iterate through all playlists
	return tracks
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
		log.Printf("Unable to shuffle and scrub playlist for playlist ID: %v", origPlaylistID)
		return nil, err
	}

	return tracks, nil
}

// This function will trim off a lot of data unnecessary to add tracks to the playlist
// The assumption is if we keep the payload small, we can add over 100 tracks to a playlist
func scrubTracks(tracks *apple.AppleResponse) *apple.AppleTrackRequest {
	// Create new Apple object to hold just IDs
	var scrubbedTracks *apple.AppleTrackRequest
	scrubbedTracks = new(apple.AppleTrackRequest)

	// Loop through returned object and pull out just the ID
	for _, track := range tracks.Data {
		var scrubbedTrack apple.TrackData
		scrubbedTrack.ID = track.ID
		scrubbedTracks.Data = append(scrubbedTracks.Data, scrubbedTrack)
	}

	return scrubbedTracks
}

// This function calculates the pagination and how many times we have to call Apple Music
func calculateOffset(trackCount int, grouping int) int {
	offsetCount := trackCount / grouping
	if trackCount%grouping > 0 {
		offsetCount++
	}
	return offsetCount
}

// Quick function to see how big the map is
func lenSyncMap(m *sync.Map) int {
	var i int
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

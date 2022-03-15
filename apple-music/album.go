package apple_music

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetAlbum(id int) Album {

	url := "https://api.music.apple.com/v1/catalog/us/albums/" + strconv.Itoa(id)

	// Create a Bearer string by appending string access token
	var secret = getSecret()
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

	var responseObject AppleAlbum
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

type AppleAlbum struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Href       string `json:"href"`
		Attributes struct {
			Artwork struct {
				Width      int    `json:"width"`
				Height     int    `json:"height"`
				URL        string `json:"url"`
				BgColor    string `json:"bgColor"`
				TextColor1 string `json:"textColor1"`
				TextColor2 string `json:"textColor2"`
				TextColor3 string `json:"textColor3"`
				TextColor4 string `json:"textColor4"`
			} `json:"artwork"`
			ArtistName          string   `json:"artistName"`
			IsSingle            bool     `json:"isSingle"`
			URL                 string   `json:"url"`
			IsComplete          bool     `json:"isComplete"`
			GenreNames          []string `json:"genreNames"`
			TrackCount          int      `json:"trackCount"`
			IsMasteredForItunes bool     `json:"isMasteredForItunes"`
			ReleaseDate         string   `json:"releaseDate"`
			Name                string   `json:"name"`
			RecordLabel         string   `json:"recordLabel"`
			Upc                 string   `json:"upc"`
			Copyright           string   `json:"copyright"`
			PlayParams          struct {
				ID   string `json:"id"`
				Kind string `json:"kind"`
			} `json:"playParams"`
			EditorialNotes struct {
				Standard string `json:"standard"`
				Short    string `json:"short"`
			} `json:"editorialNotes"`
			IsCompilation bool `json:"isCompilation"`
		} `json:"attributes"`
		Relationships struct {
			Artists struct {
				Href string `json:"href"`
				Data []struct {
					ID   string `json:"id"`
					Type string `json:"type"`
					Href string `json:"href"`
				} `json:"data"`
			} `json:"artists"`
			Tracks struct {
				Href string `json:"href"`
				Data []struct {
					ID         string `json:"id"`
					Type       string `json:"type"`
					Href       string `json:"href"`
					Attributes struct {
						Previews []struct {
							URL string `json:"url"`
						} `json:"previews"`
						Artwork struct {
							Width      int    `json:"width"`
							Height     int    `json:"height"`
							URL        string `json:"url"`
							BgColor    string `json:"bgColor"`
							TextColor1 string `json:"textColor1"`
							TextColor2 string `json:"textColor2"`
							TextColor3 string `json:"textColor3"`
							TextColor4 string `json:"textColor4"`
						} `json:"artwork"`
						ArtistName       string   `json:"artistName"`
						URL              string   `json:"url"`
						DiscNumber       int      `json:"discNumber"`
						GenreNames       []string `json:"genreNames"`
						DurationInMillis int      `json:"durationInMillis"`
						ReleaseDate      string   `json:"releaseDate"`
						Name             string   `json:"name"`
						Isrc             string   `json:"isrc"`
						HasLyrics        bool     `json:"hasLyrics"`
						AlbumName        string   `json:"albumName"`
						PlayParams       struct {
							ID   string `json:"id"`
							Kind string `json:"kind"`
						} `json:"playParams"`
						TrackNumber  int    `json:"trackNumber"`
						ComposerName string `json:"composerName"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"tracks"`
		} `json:"relationships"`
	} `json:"data"`
}

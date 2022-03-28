package main

import (
	"fuji-alexa/internal/app"
	alexa2 "fuji-alexa/internal/models/alexa"
	"github.com/aws/aws-lambda-go/lambda"
	"strconv"
	"strings"
)

func HandleFavoriteAlbumIntent(request alexa2.Request) alexa2.Response {
	//return alexa.NewSimpleResponse("Frontpage Deals", "Frontpage deal data here")
	var builder alexa2.SSMLBuilder
	builder.Say("Here are your favorite album:")
	builder.Pause("1000")
	var album app.Album
	album = app.GetAlbum(310730204)
	builder.Say(album.ArtistName)
	builder.Pause("500")
	builder.Say(album.Name)
	return alexa2.NewSSMLResponse("Favorite Album", builder.Build())
}

func HandleNumberOfPlaylistsIntent(request alexa2.Request) alexa2.Response {
	var builder alexa2.SSMLBuilder
	count := app.GetPlaylistCount(request.Session.User.UserID)
	builder.Say("You have ")
	builder.Say(strconv.Itoa(count))
	builder.Say(" playlists in your Apple Music Library")
	return alexa2.NewSSMLResponse("Playlist Count", builder.Build())
}

func HandleShuffleIntent(request alexa2.Request) alexa2.Response {
	var builder alexa2.SSMLBuilder

	playlistName := strings.ToLower(request.Body.Intent.Slots["playlistName"].Value)
	playlistID := app.FindPlaylist(request.Session.User.UserID, playlistName)

	if playlistID == "" {
		builder.Say("Apologies.  I am unable to to find the playlist " + playlistName)
		return alexa2.NewSSMLResponse("Playlist Shuffled", builder.Build())
	}

	newName, err := app.ShufflePlaylist(request.Session.User.UserID, playlistID)
	if err != nil {
		builder.Say("Apologies.  I am unable to shuffle the playlist at this moment.")
		return alexa2.NewSSMLResponse("Playlist Shuffled", builder.Build())
	}

	builder.Say("We have shuffled your playlists into a new playlist called " + newName)
	return alexa2.NewSSMLResponse("Playlist Shuffled", builder.Build())
}

func HandleHelpIntent(request alexa2.Request) alexa2.Response {
	//return alexa.NewSimpleResponse("Help", "Help regarding the available commands here")
	var builder alexa2.SSMLBuilder
	builder.Say("Here are some of the things you can ask:")
	builder.Pause("1000")
	builder.Say("What are my favorite albums.")
	builder.Pause("1000")
	builder.Say("Shuffle my playlist.")
	return alexa2.NewSSMLResponse("Fuji Music Help", builder.Build())
}

func HandleAboutIntent(request alexa2.Request) alexa2.Response {
	return alexa2.NewSimpleResponse(
		"About",
		"Fuji Music is a project created by Christopher Sheridan to enhance the Apple Music in Alexa.")
}

func IntentDispatcher(request alexa2.Request) alexa2.Response {
	var response alexa2.Response
	switch request.Body.Intent.Name {
	case "FavoriteAlbum":
		response = HandleFavoriteAlbumIntent(request)
	case "PlaylistCount":
		response = HandleNumberOfPlaylistsIntent(request)
	case "ShufflePlaylist":
		response = HandleShuffleIntent(request)
	case alexa2.HelpIntent:
		response = HandleHelpIntent(request)
	case "AboutIntent":
		response = HandleAboutIntent(request)
	default:
		response = HandleAboutIntent(request)
	}
	return response
}

type FeedResponse struct {
	Channel struct {
		Item []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

// Handler represents the Handler of lambda
func Handler(request alexa2.Request) (alexa2.Response, error) {
	return IntentDispatcher(request), nil
}

func main() {
	lambda.Start(Handler)
}

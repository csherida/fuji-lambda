package main

import (
	"fuji-alexa/alexa"
	"fuji-alexa/apple-music"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleFavoriteAlbumIntent(request alexa.Request) alexa.Response {
	//return alexa.NewSimpleResponse("Frontpage Deals", "Frontpage deal data here")
	var builder alexa.SSMLBuilder
	builder.Say("Here are your favorite album:")
	builder.Pause("1000")
	var album apple_music.Album
	album = apple_music.GetAlbum(310730204)
	builder.Say(album.ArtistName)
	builder.Pause("500")
	builder.Say(album.Name)
	return alexa.NewSSMLResponse("Favorite Album", builder.Build())
}

func HandleHelpIntent(request alexa.Request) alexa.Response {
	//return alexa.NewSimpleResponse("Help", "Help regarding the available commands here")
	var builder alexa.SSMLBuilder
	builder.Say("Here are some of the things you can ask:")
	builder.Pause("1000")
	builder.Say("Give me the frontpage deals.")
	builder.Pause("1000")
	builder.Say("Give me the popular deals.")
	return alexa.NewSSMLResponse("Slick Dealer Help", builder.Build())
}

func HandleAboutIntent(request alexa.Request) alexa.Response {
	return alexa.NewSimpleResponse("About", "Slick Dealer was created by Nic Raboy in Tracy, California as an unofficial Slick Deals application.")
}

func IntentDispatcher(request alexa.Request) alexa.Response {
	var response alexa.Response
	switch request.Body.Intent.Name {
	case "FavoriteAlbum":
		response = HandleFavoriteAlbumIntent(request)
	case alexa.HelpIntent:
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
func Handler(request alexa.Request) (alexa.Response, error) {
	return IntentDispatcher(request), nil
}

func main() {
	lambda.Start(Handler)
}

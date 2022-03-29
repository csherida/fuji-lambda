package app

import (
	models "fuji-alexa/internal/models/fuji"
	"log"
	"strconv"
)

func GetAlbum(id int) Album {

	url := "https://api.music.apple.com/v1/catalog/us/albums/" + strconv.Itoa(id)

	// AppleUserToken not needed for catalog requests
	responseObject, err := fetchAppleMusicData(&models.FujiAccount{AppleToken: ""}, url)
	if err != nil {
		log.Fatalf("Unable to get album %v", id)
		return Album{}
	}

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

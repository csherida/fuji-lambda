package models

type FujiPlaylists struct {
	FujiPlaylist         []FujiPlaylist `json:"fujiPlaylist"`
	OverallPlaylistCount int            `json:"overallPlaylistCount""`
}

type FujiPlaylist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

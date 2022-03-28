package apple

type AppleTrackRequest struct {
	Data []TrackData `json:"data"`
}

type TrackData struct {
	ID string `json:"id"`
}

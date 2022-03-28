package apple

type PlaylistRequest struct {
	Attributes struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"attributes"`
	Relationships struct {
		Parent struct {
			Data []ParentFolderData `json:"data"`
		} `json:"parent"`
	} `json:"relationships"`
}

type ParentFolderData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

package models

type FujiAccount struct {
	FujiID      string `json:"fujiID"`
	AmazonToken string `json:"amazonToken"`
	AppleToken  string `json:"appleToken"`
}

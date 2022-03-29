package app

import (
	"testing"
)

func Test_fetchAppleMusicData(t *testing.T) {
	type args struct {
		amazonToken string
		url         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Playlist Count",
			args: args{
				url:         "https://api.music.apple.com/v1/me/library/playlists",
				amazonToken: "amzn1.ask.account.testUser",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct, _ := GetFujiAccount(tt.args.amazonToken)
			got, err := fetchAppleMusicData(acct, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchAppleMusicData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Data) <= 0 {
				t.Errorf("fetchAppleMusicData() is empty")
			}
		})
	}
}

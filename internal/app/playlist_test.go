package app

import (
	"strings"
	"testing"
)

func TestGetPlaylistCount(t *testing.T) {
	type args struct {
		amazonToken string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Playlist Count",
			args: args{
				amazonToken: "amzn1.ask.account.testUser",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPlaylistCount(tt.args.amazonToken); got < tt.want {
				t.Errorf("GetPlaylistCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreatePlaylist(t *testing.T) {
	type args struct {
		amazonToken      string
		origPlaylistName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Playlist Creation Test",
			args: args{
				amazonToken:      "amzn1.ask.account.testUser",
				origPlaylistName: "Unit Testing Playlist Creation",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreatePlaylist(tt.args.amazonToken, tt.args.origPlaylistName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePlaylist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(got.Name, tt.args.origPlaylistName+" ") {
				t.Errorf("CreatePlaylist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

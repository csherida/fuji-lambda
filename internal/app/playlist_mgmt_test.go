package app

import (
	"log"
	"strings"
	"testing"
)

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

func TestGetPlaylists(t *testing.T) {
	type args struct {
		amazonToken string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "More than 25 Playlists",
			args: args{
				amazonToken: "amzn1.ask.account.testUser",
			},
			want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPlaylists(tt.args.amazonToken)
			for _, playlist := range got.FujiPlaylist {
				log.Printf("%v %v", playlist.ID, playlist.Name)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlaylists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.OverallPlaylistCount <= tt.want {
				t.Errorf("GetPlaylists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindPlaylist(t *testing.T) {
	type args struct {
		amazonToken      string
		playlistNameRqst string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Find Chill Tunes playlist", // specific to Chris Sheridan's library
			args: args{
				amazonToken:      "amzn1.ask.account.testUser",
				playlistNameRqst: "All Chill Tunes",
			},
			want: "p.5x1WhOxAz9v",
		},
		{
			name: "Negative Test", // specific to Chris Sheridan's library
			args: args{
				amazonToken:      "amzn1.ask.account.testUser",
				playlistNameRqst: "Bogus Bogus Bogus",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindPlaylist(tt.args.amazonToken, tt.args.playlistNameRqst); got != tt.want {
				t.Errorf("FindPlaylist() = %v, want %v", got, tt.want)
			}
		})
	}
}

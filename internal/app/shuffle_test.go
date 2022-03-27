package app

import (
	"fuji-alexa/internal/models/apple"
	"testing"
)

func TestShufflePlaylist(t *testing.T) {
	type args struct {
		amazonToken string
		playlistID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Shuffle Test and Write to New Playlist",
			args: args{
				amazonToken: "amzn1.ask.account.testUser",
				playlistID:  "p.oOlRRflxbK9Q",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if newName, err := ShufflePlaylist(tt.args.amazonToken, tt.args.playlistID); (err != nil) != tt.wantErr {
				t.Errorf("ShufflePlaylist() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("Unable to create new playlist %v", newName)
			}
		})
	}
}

func Test_calculateOffset(t *testing.T) {
	type args struct {
		trackCount int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Multiple offsets",
			args: args{
				trackCount: 808,
			},
			want: 9,
		},
		{
			name: "Zero tracks",
			args: args{
				trackCount: 0,
			},
			want: 0,
		},
		{
			name: "One offset",
			args: args{
				trackCount: 8,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateOffset(tt.args.trackCount); got != tt.want {
				t.Errorf("calculateOffset() = %v, doNotWant %v", got, tt.want)
			}
		})
	}
}

func Test_getTracks(t *testing.T) {
	type args struct {
		amazonToken    string
		origPlaylistID string
		pageOffset     []int
	}
	tests := []struct {
		name    string
		args    args
		want    *apple.AppleResponse
		wantErr bool
	}{
		{
			name: "Shuffle Get Tracks Test",
			args: args{
				amazonToken:    "amzn1.ask.account.testUser",
				origPlaylistID: "p.oOlRRflxbK9Q",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTracks(tt.args.amazonToken, tt.args.origPlaylistID, tt.args.pageOffset...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTracks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Data[0].ID != "" {
				t.Errorf("getTracks() got = %v, wanted empty string or nil", got)
			}
		})
	}
}

/*
func Test_shuffle(t *testing.T) {
	type args struct {
		tracks    apple.AppleTrackRequest
	}
	tests := []struct {
		name      string
		args      args
		doNotWant string
		wantErr   bool
	}{
		{
			name: "Test Shuffle",
			args: args{
				tracks: apple.AppleTrackRequest{Data: new(apple.TrackData) [apple.TrackData{ID: "abcdefg"},]}
			},
			doNotWant: "",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shuffle(tt.args.amazonToken, tt.args.origPlaylistID)
			if (err != nil) != tt.wantErr {
				t.Errorf("shuffle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.doNotWant {
				t.Errorf("shuffle() got = %v, doNotWant %v", got, tt.doNotWant)
			}
		})
	}
}
*/

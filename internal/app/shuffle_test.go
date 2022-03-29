package app

import (
	"testing"
)

func TestShufflePlaylist(t *testing.T) {
	type args struct {
		amazonToken  string
		playlistID   string
		playlistName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Shuffle Test and Write to New Playlist",
			args: args{
				amazonToken:  "amzn1.ask.account.testUser",
				playlistID:   "p.5x1WhOxAz9v", //"p.oOlRRflxbK9Q",
				playlistName: "All Chill Tunes",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct, _ := GetFujiAccount(tt.args.amazonToken)
			if newName, err := ShufflePlaylist(acct, tt.args.playlistID, tt.args.playlistName); (err != nil) != tt.wantErr {
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
			if got := calculateOffset(tt.args.trackCount, 100); got != tt.want {
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
		want    int
		wantErr bool
	}{
		{
			name: "Shuffle Get Tracks Test",
			args: args{
				amazonToken:    "amzn1.ask.account.testUser",
				origPlaylistID: "p.oOlRRflxbK9Q",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "No tracks exist for shuffle getTracks()",
			args: args{
				amazonToken:    "amzn1.ask.account.testUser",
				origPlaylistID: "p.111111111111",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct, _ := GetFujiAccount(tt.args.amazonToken)
			got, err := getTracks(acct, tt.args.origPlaylistID, tt.args.pageOffset...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTracks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && (len(got.Data) < tt.want) {
				t.Errorf("getTracks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

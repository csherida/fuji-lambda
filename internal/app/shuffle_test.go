package app

import (
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
			name: "Shuffle Test",
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
				t.Errorf("calculateOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

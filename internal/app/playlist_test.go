package app

import "testing"

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

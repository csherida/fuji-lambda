package app

import (
	"reflect"
	"testing"
)

func Test_getAlbum(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want Album
	}{
		{
			name: "Born to Run",
			args: args{
				id: 310730204,
			},
			want: Album{
				ArtistName: "Bruce Springsteen",
				Name:       "Born to Run",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAlbum(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAlbum() = %v, want %v", got, tt.want)
			}
		})
	}
}

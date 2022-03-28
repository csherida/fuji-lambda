package app

import "testing"

func Test_getSecret(t *testing.T) {
	type args struct {
		secretName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Apple Music Token",
			args: args{
				secretName: "FujiAppleMusicToken",
			},
			want: "",
		},
		{
			name: "Fuji Account Service API Key",
			args: args{
				secretName: "FujiAccountAPIKey",
			},
			want: "",
		},
		{
			name: "Apple Music Token from cache",
			args: args{
				secretName: "FujiAppleMusicToken",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSecret(tt.args.secretName); got == tt.want {
				t.Errorf("getSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

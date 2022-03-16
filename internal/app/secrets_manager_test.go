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
			name: "Born to Run",
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

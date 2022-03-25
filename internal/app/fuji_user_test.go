package app

import "testing"

func Test_getAppleUserToken(t *testing.T) {
	type args struct {
		amazonToken string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Apple Music Token",
			args: args{
				amazonToken: "amzn1.ask.account.testUser",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAppleUserToken(tt.args.amazonToken); got == tt.want {
				t.Errorf("getAppleUserToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

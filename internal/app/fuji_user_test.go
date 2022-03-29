package app

import (
	models "fuji-alexa/internal/models/fuji"
	"reflect"
	"strings"
	"testing"
)

func TestGetFujiAccount(t *testing.T) {
	type args struct {
		amazonToken string
	}
	tests := []struct {
		name       string
		args       args
		want       *models.FujiAccount
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Valid Amazon User",
			args: args{
				amazonToken: "amzn1.ask.account.testUser",
			},
			want: nil,
		},
		{
			name: "Invalid Amazon User",
			args: args{
				amazonToken: "amzn1.ask.account.bogusUser",
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFujiAccount(tt.args.amazonToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFujiAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && (tt.wantErrMsg != "") && !(strings.Contains(err.Error(), tt.wantErrMsg)) {
				t.Errorf("GetFujiAccount() error = %v, wantErrMsg %v", err, tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if got.AppleToken == "" {
					t.Errorf("GetFujiAccount() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

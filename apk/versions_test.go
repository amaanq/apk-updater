package apk

import (
	"testing"
)

func TestGetAllVersions(t *testing.T) {
	t.Run("Versions", func(t *testing.T) {
		vers, err := GetAllVersions(ClashofClans.URL)
		if err != nil {
			Log.Errorf("Versions() error = %v", err)
			return
		}
		for _, v := range vers {
			Log.Infof("%+v", v)
		}
	})
}

func TestGetVersions(t *testing.T) {
	type args struct {
		page int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{page: 1}, wantErr: false},
		{name: "2", args: args{page: 2}, wantErr: false},
		{name: "3", args: args{page: 3}, wantErr: false},
		{name: "4", args: args{page: 4}, wantErr: false},
		{name: "5", args: args{page: 5}, wantErr: false},
		{name: "6", args: args{page: 6}, wantErr: false},
		{name: "7", args: args{page: 7}, wantErr: false},
		{name: "8", args: args{page: 8}, wantErr: false},
		{name: "9", args: args{page: 9}, wantErr: false},
		{name: "10", args: args{page: 10}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetVersions(ClashofClans.URL, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("Version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

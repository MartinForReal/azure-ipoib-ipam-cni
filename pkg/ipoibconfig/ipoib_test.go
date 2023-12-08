package ipoibconfig

import (
	"net"
	"os"
	"reflect"
	"testing"
)

func TestGetIPoIBAddr(t *testing.T) {
	type args struct {
		filepath string
		macAddr  string
	}
	tests := []struct {
		name    string
		args    args
		want    *net.IPNet
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filepath: "testdata/.kvp_pool_0",
				macAddr:  "00:15:5d:33:ff:0b",
			},
			want: &net.IPNet{
				IP:   net.IPv4(172, 16, 1, 2),
				Mask: net.CIDRMask(16, 32),
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				filepath: "testdata/.kvp_pool_1",
				macAddr:  "00:15:5d:33:ff:0b",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(tt.args.filepath)
			if err != nil {
				t.Error(err)
			}
			got, err := GetIPoIBAddr(content, tt.args.macAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIPoIBAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIPoIBAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

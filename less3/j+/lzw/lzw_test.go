package main

import (
	"reflect"
	"testing"
)

func Test_LZWEncode(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want []uint16
	}{
		{
			"1",
			args{[]byte("ABABCBABABCAD")},
			[]uint16{65, 66, 256, 67, 66, 65, 258, 66, 67, 65, 68},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LZWEncode(tt.args.input, 16); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nencode() = %v, \nwant %v", got, tt.want)
			}
		})
	}
}

func Test_LZWDecode(t *testing.T) {
	type args struct {
		input []uint16
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"1",
			args{[]uint16{65, 66, 256, 67, 66, 65, 258, 66, 67, 65, 68}},
			[]byte("ABABCBABABCAD"),
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LZWDecode(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ndecode() = %s, \nwant %s", got, tt.want)
			}
		})
	}
}

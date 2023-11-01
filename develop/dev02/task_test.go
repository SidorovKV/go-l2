package main

import "testing"

func Test_unpack(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{"a4bc2d5e"}, "aaaabccddddde", false},
		{"test2", args{"abcd"}, "abcd", false},
		{"test3", args{""}, "", false},
		{"test4", args{"45"}, "", true},
		{"test5", args{"a11bc2d5e"}, "aaaaaaaaaaabccddddde", false},
		{"test6", args{"45a"}, "", true},
		{"test7", args{" 3"}, "   ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpack(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unpack() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unpackEscape(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{"a4bc2d5e"}, "aaaabccddddde", false},
		{"test2", args{"abcd"}, "abcd", false},
		{"test3", args{""}, "", false},
		{"test4", args{"45"}, "", true},
		{"test5", args{"a11bc2d5e"}, "aaaaaaaaaaabccddddde", false},
		{"test6", args{"45a"}, "", true},
		{"test7", args{" 3"}, "   ", false},
		{"test8", args{`qwe\4\5`}, "qwe45", false},
		{"test9", args{`qwe\45`}, "qwe44444", false},
		{"test10", args{`qwe\\5`}, `qwe\\\\\`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpackEscape(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackEscape() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unpackEscape() got = %v, want %v", got, tt.want)
			}
		})
	}
}

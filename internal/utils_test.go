package internal

import "testing"

func Test_IsHex(t *testing.T) {
	type args struct {
		ch byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"0", args{'0'}, true},
		{"1", args{'1'}, true},
		{"G", args{'G'}, false},
		{"g", args{'g'}, false},
		{"A", args{'A'}, true},
		{"a", args{'a'}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHex(tt.args.ch); got != tt.want {
				t.Errorf("ishex() = %v, want %v", got, tt.want)
			}
		})
	}
}

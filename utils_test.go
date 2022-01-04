package sock_pair_in_golang

import (
	"reflect"
	"testing"
)

func TestGenerateSocks_pairs(t *testing.T) {
	type args struct {
		colors        []string
		patterns      []string
		numDuplicates int
	}
	tests := []struct {
		name string
		args args
		want Socks
	}{
		{
			"3 colors, 3 patterns, 1 pair of each style",
			args{
				[]string{"red", "blue", "green"},
				[]string{"plain", "checkered", "herringbone"},
				1,
			},
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"blue", "plain", false},
				Sock{"green", "plain", true},
				Sock{"green", "plain", false},
				Sock{"red", "checkered", true},
				Sock{"red", "checkered", false},
				Sock{"blue", "checkered", true},
				Sock{"blue", "checkered", false},
				Sock{"green", "checkered", true},
				Sock{"green", "checkered", false},
				Sock{"red", "herringbone", true},
				Sock{"red", "herringbone", false},
				Sock{"blue", "herringbone", true},
				Sock{"blue", "herringbone", false},
				Sock{"green", "herringbone", true},
				Sock{"green", "herringbone", false},
			},
		},
		{
			"3 colors, 3 patterns, 0 pairs",
			args{
				[]string{"red", "blue", "green"},
				[]string{"plain", "checkered", "herringbone"},
				0,
			},
			make(Socks, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSocks(tt.args.colors, tt.args.patterns, tt.args.numDuplicates, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateSocks_singles(t *testing.T) {
	type args struct {
		colors        []string
		patterns      []string
		numDuplicates int
	}
	tests := []struct {
		name string
		args args
		want Socks
	}{
		{
			"3 colors, 3 patterns, 1 of each style",
			args{
				[]string{"red", "blue", "green"},
				[]string{"plain", "checkered", "herringbone"},
				1,
			},
			Socks{
				Sock{"red", "plain", true},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
				Sock{"red", "checkered", true},
				Sock{"blue", "checkered", true},
				Sock{"green", "checkered", true},
				Sock{"red", "herringbone", true},
				Sock{"blue", "herringbone", true},
				Sock{"green", "herringbone", true},
			},
		},
		{
			"3 colors, 3 patterns, 0 of each style",
			args{
				[]string{"red", "blue", "green"},
				[]string{"plain", "checkered", "herringbone"},
				0,
			},
			make(Socks, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSocks(tt.args.colors, tt.args.patterns, tt.args.numDuplicates, true); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

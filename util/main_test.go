package util

import (
	"reflect"
	"testing"
)

func TestSplitSpace(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "hankaku",
			args: args{
				text: "keyword1 keyword2",
			},
			want: []string{"keyword1", "keyword2"},
		},
		{
			name: "zenkaku",
			args: args{
				text: "keyword1　keyword2",
			},
			want: []string{"keyword1", "keyword2"},
		},
		{
			name: "hankaku, zenkaku mixed",
			args: args{
				text: "keyword1 　keyword2",
			},
			want: []string{"keyword1", "keyword2"},
		},
		{
			name: "hankaku renzoku",
			args: args{
				text: "keyword1  keyword2",
			},
			want: []string{"keyword1", "keyword2"},
		},
		{
			name: "first space",
			args: args{
				text: " keyword1 keyword2",
			},
			want: []string{"keyword1", "keyword2"},
		},
		{
			name: "empty text",
			args: args{
				text: "",
			},
			want: []string{},
		},
		{
			name: "space only",
			args: args{
				text: "   ",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitSpace(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

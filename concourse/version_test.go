package concourse

import (
	"math"
	"reflect"
	"testing"
)

func TestVersion_IDNumber(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		commentID string
		want      uint
	}{
		{
			name:      "empty string should return 0",
			commentID: "",
			want:      0,
		},
		{
			name:      "invalid string should return 0",
			commentID: "hello world",
			want:      0,
		},
		{
			name:      "valid string should result in 123456",
			commentID: "123456",
			want:      123456,
		},
		{
			name:      "maximum int64 should be correctly converted",
			commentID: "18446744073709551615",
			want:      math.MaxUint64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				CommentID: tt.commentID,
			}
			if got := v.IDNumber(); got != tt.want {
				t.Errorf("IDNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionFromNumber(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		commentID uint
		want      Version
	}{
		{
			name:      "empty string should return 0",
			want:      Version{CommentID: "0"},
			commentID: 0,
		},
		{
			name:      "invalid string should return 0",
			want:      Version{CommentID: "0"},
			commentID: 0,
		},
		{
			name:      "valid string should result in 123456",
			want:      Version{CommentID: "123456"},
			commentID: 123456,
		},
		{
			name:      "maximum int64 should be correctly converted",
			want:      Version{CommentID: "18446744073709551615"},
			commentID: math.MaxUint64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VersionFromNumber(tt.commentID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VersionFromNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

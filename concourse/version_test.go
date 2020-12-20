package concourse

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestVersion_MarshalJSON(t *testing.T) {

	tests := []struct {
		name      string
		commentID uint
		want      []byte
		wantErr   bool
	}{
		{
			name:      "test 0",
			commentID: 0,
			want:      bytes.NewBufferString(`{"comment_id":"0"}`).Bytes(),
			wantErr:   false,
		},
		{
			name:      "test max uint64",
			commentID: math.MaxUint64,
			want:      bytes.NewBufferString(`{"comment_id":"18446744073709551615"}`).Bytes(),
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				CommentID: tt.commentID,
			}
			got, err := v.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestVersion_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		commentID uint
		data      []byte
		wantErr   bool
	}{
		{
			name:      "no data",
			commentID: 0,
			data:      nil,
			wantErr:   false,
		},
		{
			name:      "valid comment ID",
			commentID: 123456,
			data:      bytes.NewBufferString(`{"comment_id":"123456"}"`).Bytes(),
			wantErr:   false,
		},
		{
			name:      "maximum comment ID",
			commentID: math.MaxUint64,
			data:      bytes.NewBufferString(`{"comment_id":"18446744073709551615"}"`).Bytes(),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				CommentID: tt.commentID,
			}
			if err := v.UnmarshalJSON(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package concourse

import (
	"encoding/json"
	"strconv"
)

type (
	Version struct {
		CommentID uint `json:"comment_id"`
	}

	strVersion struct {
		CommentID string `json:"comment_id"`
	}
)

// MarshalJSON marshals the uint64 to a string.
func (v *Version) MarshalJSON() ([]byte, error) {
	ver := strVersion{CommentID: strconv.FormatUint(uint64(v.CommentID), 10)}

	return json.Marshal(&ver)
}

// UnmarshalJSON never fails. On failure it sets the CommentID to 0.
func (v *Version) UnmarshalJSON(b []byte) error {
	var ver strVersion
	err := json.Unmarshal(b, &ver)
	if err != nil {
		v.CommentID = 0
		return nil
	}

	// The error is ignored because the signal of failure is CommentID = 0.
	id, _ := strconv.Atoi(ver.CommentID)
	v.CommentID = uint(id)
	return nil
}

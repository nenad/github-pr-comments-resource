package concourse

import (
	"strconv"
)

type (
	// Version represents a single version of a resource.
	Version struct {
		// CommentID is a unique identifier to a GitHub comment on a pull request.
		CommentID string `json:"comment_id"`
	}
)

// VersionFromNumber creates a Version from a uint instead of a string.
func VersionFromNumber(commentID uint) Version {
	return Version{CommentID: strconv.FormatUint(uint64(commentID), 10)}
}

// IDNumber returns the version identifier in a numeric value.
func (v *Version) IDNumber() uint {
	num, _ := strconv.ParseUint(v.CommentID, 10, 64)

	return uint(num)
}

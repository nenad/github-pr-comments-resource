package github

import (
	"regexp"
)

// FilterSince returns comments newer than the one given.
func FilterSince(commentID uint) Filter {
	return func(comments []Comment) []Comment {
		var newComments []Comment
		for _, c := range comments {
			if c.ID < commentID {
				continue
			}
			newComments = append(newComments, c)
		}

		return newComments
	}
}

// FilterCommentRegex returns comments that match a given regex.
func FilterCommentRegex(regex string) Filter {
	return func(comments []Comment) []Comment {
		matcher := regexp.MustCompile(regex)
		var newComments []Comment
		for _, c := range comments {
			if !matcher.MatchString(c.Body) {
				continue
			}
			newComments = append(newComments, c)
		}

		return newComments
	}
}

// FilterLabels returns comment that match at least one of the given labels.
func FilterLabels(labels []string) Filter {
	return func(comments []Comment) []Comment {
		var newComments []Comment
		for _, c := range comments {
			for _, l := range labels {
				if c.HasLabel(l) {
					newComments = append(newComments, c)
					break
				}
			}
		}

		return newComments
	}
}

// FilterLatestPerPR returns the latest comment per PR that matches previous criteria if any.
func FilterLatestPerPR() Filter {
	return func(comments []Comment) []Comment {
		prMap := make(map[uint]Comment)
		for _, c := range comments {
			val, ok := prMap[c.PullRequestNumber]
			if !ok || c.ID > val.ID {
				prMap[c.PullRequestNumber] = c
			}
		}

		var newComments []Comment
		for _, c := range prMap {
			newComments = append(newComments, c)
		}

		return newComments
	}
}

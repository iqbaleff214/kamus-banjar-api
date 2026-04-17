package community

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrWordNotFound   = errors.New("word not found or not active")
	ErrCommentNotFound = errors.New("comment not found")
	ErrForbidden      = errors.New("insufficient permissions")
	ErrInvalidVote    = errors.New("vote must be 1 (upvote) or -1 (downvote)")
	ErrEmptyBody      = errors.New("comment body cannot be empty")
	ErrNoWordToday    = errors.New("no active word available for today")
)

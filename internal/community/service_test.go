package community_test

import (
	"errors"
	"testing"

	"github.com/iqbaleff214/kamus-banjar-api/internal/community"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newService(repo community.Repository) community.Service {
	return community.NewService(repo)
}

// ─────────────────────────────────────────────────────────────
// Votes
// ─────────────────────────────────────────────────────────────

func TestCastVote_UpvoteNew(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	repo.On("GetUserVote", "abah", "u1").Return(0, nil)
	repo.On("UpsertVote", "abah", "u1", 1).Return(nil)

	err := newService(repo).CastVote("abah", "u1", 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCastVote_Toggle_RemovesVote(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	repo.On("GetUserVote", "abah", "u1").Return(1, nil) // already upvoted
	repo.On("DeleteVote", "abah", "u1").Return(nil)

	err := newService(repo).CastVote("abah", "u1", 1) // same vote → toggle off

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCastVote_ChangeDirection(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	repo.On("GetUserVote", "abah", "u1").Return(1, nil) // was up
	repo.On("UpsertVote", "abah", "u1", -1).Return(nil) // now down

	err := newService(repo).CastVote("abah", "u1", -1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCastVote_InvalidValue(t *testing.T) {
	err := newService(new(community.MockRepository)).CastVote("abah", "u1", 0)
	assert.ErrorIs(t, err, community.ErrInvalidVote)
}

func TestCastVote_WordNotFound(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "ghost").Return(false, nil)

	err := newService(repo).CastVote("ghost", "u1", 1)

	assert.ErrorIs(t, err, community.ErrWordNotFound)
}

func TestGetVoteSummary(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("GetVoteSummary", "abah").Return(5, 2, nil)

	up, down, err := newService(repo).GetVoteSummary("abah")

	assert.NoError(t, err)
	assert.Equal(t, 5, up)
	assert.Equal(t, 2, down)
}

// ─────────────────────────────────────────────────────────────
// Comments
// ─────────────────────────────────────────────────────────────

func TestPostComment_Success(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	repo.On("CreateComment", mock.AnythingOfType("community.Comment")).Return(nil)

	c, err := newService(repo).PostComment("abah", "u1", "Alice", "Great word!", nil)

	assert.NoError(t, err)
	assert.Equal(t, "Great word!", c.Body)
	assert.Equal(t, "u1", c.UserID)
	assert.NotEmpty(t, c.ID)
	repo.AssertExpectations(t)
}

func TestPostComment_EmptyBody(t *testing.T) {
	_, err := newService(new(community.MockRepository)).PostComment("abah", "u1", "Alice", "   ", nil)
	assert.ErrorIs(t, err, community.ErrEmptyBody)
}

func TestPostComment_WordNotFound(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "ghost").Return(false, nil)

	_, err := newService(repo).PostComment("ghost", "u1", "Alice", "hello", nil)

	assert.ErrorIs(t, err, community.ErrWordNotFound)
}

func TestPostComment_Reply_ParentNotFound(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	parentID := "bad-parent"
	repo.On("GetComment", parentID).Return(community.Comment{}, community.ErrCommentNotFound)

	_, err := newService(repo).PostComment("abah", "u1", "Alice", "reply", &parentID)

	assert.ErrorIs(t, err, community.ErrCommentNotFound)
}

func TestPostComment_Reply_WrongWord(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	parentID := "p1"
	// Parent comment belongs to a different word
	repo.On("GetComment", parentID).Return(community.Comment{ID: "p1", Word: "other"}, nil)

	_, err := newService(repo).PostComment("abah", "u1", "Alice", "reply", &parentID)

	assert.ErrorIs(t, err, community.ErrCommentNotFound)
}

func TestListComments(t *testing.T) {
	repo := new(community.MockRepository)
	comments := []community.Comment{
		{ID: "c1", Word: "abah", UserID: "u1", Body: "nice"},
	}
	repo.On("ListComments", "abah").Return(comments, nil)

	result, err := newService(repo).ListComments("abah")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDeleteComment_Owner(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("GetComment", "c1").Return(community.Comment{ID: "c1", UserID: "u1"}, nil)
	repo.On("DeleteComment", "c1").Return(nil)

	err := newService(repo).DeleteComment("c1", "u1", false)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteComment_Admin_AnyComment(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("GetComment", "c1").Return(community.Comment{ID: "c1", UserID: "u2"}, nil)
	repo.On("DeleteComment", "c1").Return(nil)

	err := newService(repo).DeleteComment("c1", "admin1", true)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteComment_NonOwner_Forbidden(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("GetComment", "c1").Return(community.Comment{ID: "c1", UserID: "u2"}, nil)

	err := newService(repo).DeleteComment("c1", "u1", false)

	assert.ErrorIs(t, err, community.ErrForbidden)
}

// ─────────────────────────────────────────────────────────────
// Bookmarks
// ─────────────────────────────────────────────────────────────

func TestAddBookmark_Success(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	repo.On("AddBookmark", "abah", "u1").Return(nil)

	err := newService(repo).AddBookmark("abah", "u1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAddBookmark_WordNotFound(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "ghost").Return(false, nil)

	err := newService(repo).AddBookmark("ghost", "u1")

	assert.ErrorIs(t, err, community.ErrWordNotFound)
}

func TestRemoveBookmark_Success(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("RemoveBookmark", "abah", "u1").Return(nil)

	err := newService(repo).RemoveBookmark("abah", "u1")

	assert.NoError(t, err)
}

func TestListBookmarks(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ListBookmarks", "u1", 1, 20).Return([]string{"abah", "abadan"}, 2, nil)

	words, total, err := newService(repo).ListBookmarks("u1", 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, []string{"abah", "abadan"}, words)
}

// ─────────────────────────────────────────────────────────────
// Word of the Day
// ─────────────────────────────────────────────────────────────

func TestGetWordOfTheDay_FromDB(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("GetWordOfTheDay", mock.AnythingOfType("string")).Return("abah", nil)

	word, err := newService(repo).GetWordOfTheDay()

	assert.NoError(t, err)
	assert.Equal(t, "abah", word)
}

func TestGetWordOfTheDay_FallsBackToRandom(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("GetWordOfTheDay", mock.AnythingOfType("string")).Return("", community.ErrNotFound)
	repo.On("GetRandomActiveWord").Return("abadan", nil)
	repo.On("SetWordOfTheDay", "abadan", "", mock.AnythingOfType("string")).Return(nil)

	word, err := newService(repo).GetWordOfTheDay()

	assert.NoError(t, err)
	assert.Equal(t, "abadan", word)
	repo.AssertExpectations(t)
}

func TestGetWordOfTheDay_NoActiveWords(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("GetWordOfTheDay", mock.AnythingOfType("string")).Return("", community.ErrNotFound)
	repo.On("GetRandomActiveWord").Return("", community.ErrNoWordToday)

	_, err := newService(repo).GetWordOfTheDay()

	assert.ErrorIs(t, err, community.ErrNoWordToday)
}

func TestSetWordOfTheDay_Success(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	repo.On("SetWordOfTheDay", "abah", "admin1", "2026-04-20").Return(nil)

	err := newService(repo).SetWordOfTheDay("abah", "admin1", "2026-04-20")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSetWordOfTheDay_InvalidDate(t *testing.T) {
	err := newService(new(community.MockRepository)).SetWordOfTheDay("abah", "admin1", "not-a-date")
	assert.Error(t, err)
}

func TestSetWordOfTheDay_WordNotFound(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "ghost").Return(false, nil)

	err := newService(repo).SetWordOfTheDay("ghost", "admin1", "2026-04-20")

	assert.ErrorIs(t, err, community.ErrWordNotFound)
}

func TestSetWordOfTheDay_DBError(t *testing.T) {
	repo := new(community.MockRepository)
	repo.On("ActiveWordExists", "abah").Return(true, nil)
	repo.On("SetWordOfTheDay", "abah", "admin1", "2026-04-20").Return(errors.New("db error"))

	err := newService(repo).SetWordOfTheDay("abah", "admin1", "2026-04-20")

	assert.Error(t, err)
}

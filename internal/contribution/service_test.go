package contribution_test

import (
	"testing"
	"time"

	"github.com/iqbaleff214/kamus-banjar-api/internal/contribution"
	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newService(repo contribution.Repository) contribution.Service {
	return contribution.NewService(repo)
}

func validReq() contribution.SubmitRequest {
	return contribution.SubmitRequest{
		Word:     "abah",
		Syllable: "a-bah",
		Alphabet: "a",
		Meanings: []dictionary.WordMeaning{
			{Definitions: []dictionary.WordDefinition{{PartOfSpeech: "n", Definition: "ayah"}}},
		},
	}
}

func pendingWord() dictionary.Word {
	return dictionary.Word{ID: 1, Word: "abah", Alphabet: "a", Status: "pending"}
}

func activeWord() dictionary.Word {
	return dictionary.Word{ID: 1, Word: "abah", Alphabet: "a", Status: "active"}
}

// ─────────────────────────────────────────────────────────────
// Submit
// ─────────────────────────────────────────────────────────────

func TestSubmit_Success(t *testing.T) {
	repo := new(contribution.MockRepository)
	repo.On("WordExists", "abah").Return(false, nil)
	repo.On("CreateWord", mock.AnythingOfType("dictionary.Word"), "u1", "community", "pending").Return(int64(1), nil)
	repo.On("LogAction", mock.AnythingOfType("contribution.Contribution")).Return(nil)

	svc := newService(repo)
	c, err := svc.Submit("u1", validReq())

	assert.NoError(t, err)
	assert.Equal(t, "submitted", c.Action)
	assert.Equal(t, int64(1), c.WordID)
	assert.Equal(t, "u1", c.ContributorID)
	assert.NotEmpty(t, c.ID)
	repo.AssertExpectations(t)
}

func TestSubmit_DuplicateWord(t *testing.T) {
	repo := new(contribution.MockRepository)
	repo.On("WordExists", "abah").Return(true, nil)

	svc := newService(repo)
	_, err := svc.Submit("u1", validReq())

	assert.ErrorIs(t, err, contribution.ErrDuplicateWord)
}

func TestSubmit_MissingWord(t *testing.T) {
	req := validReq()
	req.Word = ""
	svc := newService(new(contribution.MockRepository))
	_, err := svc.Submit("u1", req)
	assert.Error(t, err)
	assert.True(t, contribution.IsValidationErr(err))
}

func TestSubmit_MissingMeanings(t *testing.T) {
	req := validReq()
	req.Meanings = nil
	svc := newService(new(contribution.MockRepository))
	_, err := svc.Submit("u1", req)
	assert.Error(t, err)
	assert.True(t, contribution.IsValidationErr(err))
}

func TestSubmit_InvalidAlphabet(t *testing.T) {
	req := validReq()
	req.Alphabet = "ab"
	svc := newService(new(contribution.MockRepository))
	_, err := svc.Submit("u1", req)
	assert.Error(t, err)
	assert.True(t, contribution.IsValidationErr(err))
}

// ─────────────────────────────────────────────────────────────
// Edit
// ─────────────────────────────────────────────────────────────

func TestEdit_OwnPendingWord(t *testing.T) {
	repo := new(contribution.MockRepository)
	existing := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "submitted"}
	repo.On("GetByID", "c1").Return(existing, nil)
	repo.On("GetWordByID", int64(1)).Return(pendingWord(), "pending", nil)
	repo.On("UpdateWord", int64(1), mock.AnythingOfType("dictionary.Word")).Return(nil)
	repo.On("SetWordStatus", int64(1), "pending", (*string)(nil), (*time.Time)(nil)).Return(nil)
	repo.On("LogAction", mock.AnythingOfType("contribution.Contribution")).Return(nil)

	svc := newService(repo)
	c, err := svc.Edit("u1", "c1", validReq())

	assert.NoError(t, err)
	assert.Equal(t, "revised", c.Action)
	repo.AssertExpectations(t)
}

func TestEdit_OwnRejectedWord(t *testing.T) {
	repo := new(contribution.MockRepository)
	existing := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "rejected"}
	rejectedWord := dictionary.Word{ID: 1, Word: "abah", Alphabet: "a", Status: "rejected"}
	repo.On("GetByID", "c1").Return(existing, nil)
	repo.On("GetWordByID", int64(1)).Return(rejectedWord, "rejected", nil)
	repo.On("UpdateWord", int64(1), mock.AnythingOfType("dictionary.Word")).Return(nil)
	repo.On("SetWordStatus", int64(1), "pending", (*string)(nil), (*time.Time)(nil)).Return(nil)
	repo.On("LogAction", mock.AnythingOfType("contribution.Contribution")).Return(nil)

	svc := newService(repo)
	c, err := svc.Edit("u1", "c1", validReq())

	assert.NoError(t, err)
	assert.Equal(t, "revised", c.Action)
	repo.AssertExpectations(t)
}

func TestEdit_OtherUsersWord_Forbidden(t *testing.T) {
	repo := new(contribution.MockRepository)
	existing := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u2", Action: "submitted"}
	repo.On("GetByID", "c1").Return(existing, nil)

	svc := newService(repo)
	_, err := svc.Edit("u1", "c1", validReq())

	assert.ErrorIs(t, err, contribution.ErrForbidden)
}

func TestEdit_ActiveWord_Forbidden(t *testing.T) {
	repo := new(contribution.MockRepository)
	existing := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "approved"}
	repo.On("GetByID", "c1").Return(existing, nil)
	repo.On("GetWordByID", int64(1)).Return(activeWord(), "active", nil)

	svc := newService(repo)
	_, err := svc.Edit("u1", "c1", validReq())

	assert.ErrorIs(t, err, contribution.ErrWordActive)
}

// ─────────────────────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────────────────────

func TestDelete_PendingWord_Success(t *testing.T) {
	repo := new(contribution.MockRepository)
	existing := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "submitted"}
	repo.On("GetByID", "c1").Return(existing, nil)
	repo.On("GetWordByID", int64(1)).Return(pendingWord(), "pending", nil)
	repo.On("DeleteWord", int64(1)).Return(nil)

	svc := newService(repo)
	err := svc.Delete("u1", "c1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDelete_ActiveWord_Forbidden(t *testing.T) {
	repo := new(contribution.MockRepository)
	existing := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "approved"}
	repo.On("GetByID", "c1").Return(existing, nil)
	repo.On("GetWordByID", int64(1)).Return(activeWord(), "active", nil)

	svc := newService(repo)
	err := svc.Delete("u1", "c1")

	assert.ErrorIs(t, err, contribution.ErrNotPending)
}

func TestDelete_OtherUsersWord_Forbidden(t *testing.T) {
	repo := new(contribution.MockRepository)
	existing := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u2", Action: "submitted"}
	repo.On("GetByID", "c1").Return(existing, nil)

	svc := newService(repo)
	err := svc.Delete("u1", "c1")

	assert.ErrorIs(t, err, contribution.ErrForbidden)
}

// ─────────────────────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────────────────────

func TestGetByID_Owner(t *testing.T) {
	repo := new(contribution.MockRepository)
	c := contribution.Contribution{ID: "c1", ContributorID: "u1", Action: "submitted"}
	repo.On("GetByID", "c1").Return(c, nil)

	svc := newService(repo)
	result, err := svc.GetByID("u1", "c1", false)

	assert.NoError(t, err)
	assert.Equal(t, "c1", result.ID)
}

func TestGetByID_Admin_CanSeeAnyone(t *testing.T) {
	repo := new(contribution.MockRepository)
	c := contribution.Contribution{ID: "c1", ContributorID: "u2", Action: "submitted"}
	repo.On("GetByID", "c1").Return(c, nil)

	svc := newService(repo)
	result, err := svc.GetByID("admin1", "c1", true)

	assert.NoError(t, err)
	assert.Equal(t, "c1", result.ID)
}

func TestGetByID_NonOwner_Forbidden(t *testing.T) {
	repo := new(contribution.MockRepository)
	c := contribution.Contribution{ID: "c1", ContributorID: "u2", Action: "submitted"}
	repo.On("GetByID", "c1").Return(c, nil)

	svc := newService(repo)
	_, err := svc.GetByID("u1", "c1", false)

	assert.ErrorIs(t, err, contribution.ErrForbidden)
}

// ─────────────────────────────────────────────────────────────
// Approve
// ─────────────────────────────────────────────────────────────

func TestApprove_Success(t *testing.T) {
	repo := new(contribution.MockRepository)
	c := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "submitted"}
	repo.On("GetByID", "c1").Return(c, nil)
	repo.On("SetWordStatus", int64(1), "active", mock.AnythingOfType("*string"), mock.AnythingOfType("*time.Time")).Return(nil)
	repo.On("LogAction", mock.AnythingOfType("contribution.Contribution")).Return(nil)

	svc := newService(repo)
	err := svc.Approve("admin1", "c1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestApprove_NotFound(t *testing.T) {
	repo := new(contribution.MockRepository)
	repo.On("GetByID", "c1").Return(contribution.Contribution{}, contribution.ErrNotFound)

	svc := newService(repo)
	err := svc.Approve("admin1", "c1")

	assert.ErrorIs(t, err, contribution.ErrNotFound)
}

// ─────────────────────────────────────────────────────────────
// Reject
// ─────────────────────────────────────────────────────────────

func TestReject_Success(t *testing.T) {
	repo := new(contribution.MockRepository)
	c := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "submitted"}
	repo.On("GetByID", "c1").Return(c, nil)
	repo.On("SetWordStatus", int64(1), "rejected", mock.AnythingOfType("*string"), (*time.Time)(nil)).Return(nil)
	repo.On("LogAction", mock.AnythingOfType("contribution.Contribution")).Return(nil)

	svc := newService(repo)
	err := svc.Reject("admin1", "c1", "Spelling conflict with existing entry.")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestReject_WithoutNotes(t *testing.T) {
	repo := new(contribution.MockRepository)
	c := contribution.Contribution{ID: "c1", WordID: 1, ContributorID: "u1", Action: "submitted"}
	repo.On("GetByID", "c1").Return(c, nil)
	repo.On("SetWordStatus", int64(1), "rejected", mock.AnythingOfType("*string"), (*time.Time)(nil)).Return(nil)
	repo.On("LogAction", mock.AnythingOfType("contribution.Contribution")).Return(nil)

	svc := newService(repo)
	err := svc.Reject("admin1", "c1", "")

	assert.NoError(t, err)
}

// ─────────────────────────────────────────────────────────────
// AdminCreateWord
// ─────────────────────────────────────────────────────────────

func TestAdminCreateWord_Success(t *testing.T) {
	repo := new(contribution.MockRepository)
	repo.On("WordExists", "abah").Return(false, nil)
	repo.On("CreateWord", mock.AnythingOfType("dictionary.Word"), "", "official", "active").Return(int64(5), nil)

	svc := newService(repo)
	w, err := svc.AdminCreateWord(validReq())

	assert.NoError(t, err)
	assert.Equal(t, int64(5), w.ID)
	assert.Equal(t, "official", w.Source)
	assert.Equal(t, "active", w.Status)
}

func TestAdminCreateWord_Duplicate(t *testing.T) {
	repo := new(contribution.MockRepository)
	repo.On("WordExists", "abah").Return(true, nil)

	svc := newService(repo)
	_, err := svc.AdminCreateWord(validReq())

	assert.ErrorIs(t, err, contribution.ErrDuplicateWord)
}

// ─────────────────────────────────────────────────────────────
// AdminDeleteWord (soft-delete)
// ─────────────────────────────────────────────────────────────

func TestAdminDeleteWord_Success(t *testing.T) {
	repo := new(contribution.MockRepository)
	repo.On("GetWordByID", int64(1)).Return(activeWord(), "active", nil)
	repo.On("SetWordStatus", int64(1), "rejected", (*string)(nil), (*time.Time)(nil)).Return(nil)

	svc := newService(repo)
	err := svc.AdminDeleteWord(1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAdminDeleteWord_NotFound(t *testing.T) {
	repo := new(contribution.MockRepository)
	repo.On("GetWordByID", int64(99)).Return(dictionary.Word{}, "", contribution.ErrWordNotFound)

	svc := newService(repo)
	err := svc.AdminDeleteWord(99)

	assert.ErrorIs(t, err, contribution.ErrWordNotFound)
}

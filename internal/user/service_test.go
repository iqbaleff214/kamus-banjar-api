package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/iqbaleff214/kamus-banjar-api/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func mustHash(password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(h)
}

const testSecret = "super-secret-key-for-testing-only-32ch"

func newService(repo user.Repository) user.Service {
	return user.NewService(repo, testSecret, 15*time.Minute, 7*24*time.Hour)
}

// ─────────────────────────────────────────────────────────────
// Register
// ─────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("FindByEmail", "alice@example.com").Return(user.User{}, "", user.ErrNotFound)
	repo.On("Create", mock.AnythingOfType("user.User"), mock.AnythingOfType("string")).Return(nil)

	svc := newService(repo)
	u, err := svc.Register(user.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "password123"})

	assert.NoError(t, err)
	assert.Equal(t, "Alice", u.Name)
	assert.Equal(t, "alice@example.com", u.Email)
	assert.Equal(t, "user", u.Role)
	assert.True(t, u.IsActive)
	repo.AssertExpectations(t)
}

func TestRegister_EmailTaken(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("FindByEmail", "alice@example.com").Return(user.User{ID: "1"}, "hash", nil)

	svc := newService(repo)
	_, err := svc.Register(user.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "password123"})

	assert.ErrorIs(t, err, user.ErrEmailTaken)
	repo.AssertExpectations(t)
}

func TestRegister_InvalidEmail(t *testing.T) {
	svc := newService(new(user.MockRepository))
	_, err := svc.Register(user.RegisterRequest{Name: "Alice", Email: "not-an-email", Password: "password123"})
	assert.Error(t, err)
}

func TestRegister_ShortPassword(t *testing.T) {
	svc := newService(new(user.MockRepository))
	_, err := svc.Register(user.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "short"})
	assert.Error(t, err)
}

func TestRegister_ShortName(t *testing.T) {
	svc := newService(new(user.MockRepository))
	_, err := svc.Register(user.RegisterRequest{Name: "A", Email: "alice@example.com", Password: "password123"})
	assert.Error(t, err)
}

// ─────────────────────────────────────────────────────────────
// Login
// ─────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	repo := new(user.MockRepository)
	hash := mustHash("password123")
	activeUser := user.User{ID: "u1", Name: "Alice", Email: "alice@example.com", Role: "user", IsActive: true}
	repo.On("FindByEmail", "alice@example.com").Return(activeUser, hash, nil)
	repo.On("SaveRefreshToken", "u1", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	svc := newService(repo)
	pair, err := svc.Login(user.LoginRequest{Email: "alice@example.com", Password: "password123"})

	assert.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	repo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("FindByEmail", "ghost@example.com").Return(user.User{}, "", user.ErrNotFound)

	svc := newService(repo)
	_, err := svc.Login(user.LoginRequest{Email: "ghost@example.com", Password: "password123"})

	assert.ErrorIs(t, err, user.ErrInvalidCreds)
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := new(user.MockRepository)
	hash := mustHash("password123")
	activeUser := user.User{ID: "u1", IsActive: true}
	repo.On("FindByEmail", "alice@example.com").Return(activeUser, hash, nil)

	svc := newService(repo)
	_, err := svc.Login(user.LoginRequest{Email: "alice@example.com", Password: "wrongpassword"})

	assert.ErrorIs(t, err, user.ErrInvalidCreds)
}

func TestLogin_InactiveAccount(t *testing.T) {
	repo := new(user.MockRepository)
	hash := mustHash("password123")
	inactiveUser := user.User{ID: "u1", IsActive: false}
	repo.On("FindByEmail", "alice@example.com").Return(inactiveUser, hash, nil)

	svc := newService(repo)
	_, err := svc.Login(user.LoginRequest{Email: "alice@example.com", Password: "password123"})

	assert.ErrorIs(t, err, user.ErrAccountInactive)
}

func TestLogin_EmptyCredentials(t *testing.T) {
	svc := newService(new(user.MockRepository))
	_, err := svc.Login(user.LoginRequest{})
	assert.Error(t, err)
}

// ─────────────────────────────────────────────────────────────
// RefreshTokens
// ─────────────────────────────────────────────────────────────

func TestRefreshTokens_Success(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("FindRefreshToken", mock.AnythingOfType("string")).Return("u1", nil)
	repo.On("RevokeRefreshToken", mock.AnythingOfType("string")).Return(nil)
	repo.On("FindByID", "u1").Return(user.User{ID: "u1", Role: "user"}, nil)
	repo.On("SaveRefreshToken", "u1", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	svc := newService(repo)
	pair, err := svc.RefreshTokens("some-valid-refresh-token")

	assert.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	repo.AssertExpectations(t)
}

func TestRefreshTokens_InvalidToken(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("FindRefreshToken", mock.AnythingOfType("string")).Return("", user.ErrTokenInvalid)

	svc := newService(repo)
	_, err := svc.RefreshTokens("bad-token")

	assert.ErrorIs(t, err, user.ErrTokenInvalid)
}

func TestRefreshTokens_EmptyToken(t *testing.T) {
	svc := newService(new(user.MockRepository))
	_, err := svc.RefreshTokens("")
	assert.ErrorIs(t, err, user.ErrTokenInvalid)
}

// ─────────────────────────────────────────────────────────────
// Logout
// ─────────────────────────────────────────────────────────────

func TestLogout_Success(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("RevokeRefreshToken", mock.AnythingOfType("string")).Return(nil)

	svc := newService(repo)
	err := svc.Logout("some-token")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLogout_EmptyTokenIsNoop(t *testing.T) {
	svc := newService(new(user.MockRepository))
	err := svc.Logout("")
	assert.NoError(t, err)
}

// ─────────────────────────────────────────────────────────────
// Me
// ─────────────────────────────────────────────────────────────

func TestMe_Success(t *testing.T) {
	repo := new(user.MockRepository)
	expected := user.User{ID: "u1", Name: "Alice", Email: "alice@example.com"}
	repo.On("FindByID", "u1").Return(expected, nil)

	svc := newService(repo)
	u, err := svc.Me("u1")

	assert.NoError(t, err)
	assert.Equal(t, expected, u)
}

func TestMe_NotFound(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("FindByID", "missing").Return(user.User{}, user.ErrNotFound)

	svc := newService(repo)
	_, err := svc.Me("missing")

	assert.ErrorIs(t, err, user.ErrNotFound)
}

// ─────────────────────────────────────────────────────────────
// UpdateProfile
// ─────────────────────────────────────────────────────────────

func TestUpdateProfile_NameOnly(t *testing.T) {
	repo := new(user.MockRepository)
	existing := user.User{ID: "u1", Name: "Alice", Email: "alice@example.com"}
	repo.On("FindByID", "u1").Return(existing, nil)
	repo.On("Update", "u1", "Bob").Return(nil)

	svc := newService(repo)
	u, err := svc.UpdateProfile("u1", user.UpdateProfileRequest{Name: "Bob"})

	assert.NoError(t, err)
	assert.Equal(t, "Bob", u.Name)
	repo.AssertExpectations(t)
}

func TestUpdateProfile_PasswordChange(t *testing.T) {
	repo := new(user.MockRepository)
	existing := user.User{ID: "u1", Name: "Alice", Email: "alice@example.com"}
	hash := mustHash("password123")
	repo.On("FindByID", "u1").Return(existing, nil)
	repo.On("FindByEmail", "alice@example.com").Return(existing, hash, nil)
	repo.On("UpdatePassword", "u1", mock.AnythingOfType("string")).Return(nil)

	svc := newService(repo)
	_, err := svc.UpdateProfile("u1", user.UpdateProfileRequest{
		OldPassword: "password123",
		NewPassword: "newpassword123",
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateProfile_WrongOldPassword(t *testing.T) {
	repo := new(user.MockRepository)
	existing := user.User{ID: "u1", Name: "Alice", Email: "alice@example.com"}
	hash := mustHash("password123")
	repo.On("FindByID", "u1").Return(existing, nil)
	repo.On("FindByEmail", "alice@example.com").Return(existing, hash, nil)

	svc := newService(repo)
	_, err := svc.UpdateProfile("u1", user.UpdateProfileRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	})

	assert.ErrorIs(t, err, user.ErrWrongPassword)
}

func TestUpdateProfile_NewPasswordTooShort(t *testing.T) {
	repo := new(user.MockRepository)
	existing := user.User{ID: "u1", Name: "Alice", Email: "alice@example.com"}
	repo.On("FindByID", "u1").Return(existing, nil)

	svc := newService(repo)
	_, err := svc.UpdateProfile("u1", user.UpdateProfileRequest{
		OldPassword: "password123",
		NewPassword: "short",
	})

	assert.Error(t, err)
}

func TestUpdateProfile_MissingOldPassword(t *testing.T) {
	repo := new(user.MockRepository)
	existing := user.User{ID: "u1", Name: "Alice", Email: "alice@example.com"}
	repo.On("FindByID", "u1").Return(existing, nil)

	svc := newService(repo)
	_, err := svc.UpdateProfile("u1", user.UpdateProfileRequest{
		NewPassword: "newpassword123",
	})

	assert.Error(t, err)
}

// ─────────────────────────────────────────────────────────────
// ListUsers
// ─────────────────────────────────────────────────────────────

func TestListUsers(t *testing.T) {
	repo := new(user.MockRepository)
	users := []user.User{{ID: "u1"}, {ID: "u2"}}
	active := true
	repo.On("ListUsers", 1, 20, "user", &active).Return(users, 2, nil)

	svc := newService(repo)
	result, total, err := svc.ListUsers(1, 20, "user", &active)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
}

// ─────────────────────────────────────────────────────────────
// SetActive
// ─────────────────────────────────────────────────────────────

func TestSetActive_Success(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("SetActive", "u2", false).Return(nil)

	svc := newService(repo)
	err := svc.SetActive("u2", "u1", false)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSetActive_SelfModify(t *testing.T) {
	svc := newService(new(user.MockRepository))
	err := svc.SetActive("u1", "u1", false)
	assert.ErrorIs(t, err, user.ErrSelfModify)
}

// ─────────────────────────────────────────────────────────────
// Promote
// ─────────────────────────────────────────────────────────────

func TestPromote_Success(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("SetRole", "u1", "admin").Return(nil)

	svc := newService(repo)
	err := svc.Promote("u1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPromote_NotFound(t *testing.T) {
	repo := new(user.MockRepository)
	repo.On("SetRole", "missing", "admin").Return(user.ErrNotFound)

	svc := newService(repo)
	err := svc.Promote("missing")

	assert.True(t, errors.Is(err, user.ErrNotFound))
}

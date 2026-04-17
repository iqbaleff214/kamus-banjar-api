package user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service defines all business logic operations for the user domain.
type Service interface {
	Register(req RegisterRequest) (User, error)
	Login(req LoginRequest) (TokenPair, error)
	RefreshTokens(refreshToken string) (TokenPair, error)
	Logout(refreshToken string) error
	Me(userID string) (User, error)
	UpdateProfile(userID string, req UpdateProfileRequest) (User, error)
	// admin
	ListUsers(page, limit int, role string, active *bool) ([]User, int, error)
	SetActive(userID, callerID string, active bool) error
	Promote(userID string) error
}

type service struct {
	repo          Repository
	jwtSecret     string
	jwtAccessTTL  time.Duration
	jwtRefreshTTL time.Duration
}

// NewService creates a Service with the given repository and JWT config.
func NewService(repo Repository, jwtSecret string, accessTTL, refreshTTL time.Duration) Service {
	return &service{
		repo:          repo,
		jwtSecret:     jwtSecret,
		jwtAccessTTL:  accessTTL,
		jwtRefreshTTL: refreshTTL,
	}
}

// ─────────────────────────────────────────────────────────────
// Public methods
// ─────────────────────────────────────────────────────────────

func (s *service) Register(req RegisterRequest) (User, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if len(req.Name) < 2 || len(req.Name) > 100 {
		return User{}, errors.New("name must be between 2 and 100 characters")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return User{}, errors.New("invalid email address")
	}
	if len(req.Password) < 8 {
		return User{}, errors.New("password must be at least 8 characters")
	}

	// Check uniqueness
	_, _, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return User{}, ErrEmailTaken
	}
	if !errors.Is(err, ErrNotFound) {
		return User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return User{}, err
	}

	u := User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Role:     "user",
		IsActive: true,
	}

	if err = s.repo.Create(u, string(hash)); err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *service) Login(req LoginRequest) (TokenPair, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Email == "" || req.Password == "" {
		return TokenPair{}, errors.New("email and password are required")
	}

	u, hash, err := s.repo.FindByEmail(req.Email)
	if errors.Is(err, ErrNotFound) {
		return TokenPair{}, ErrInvalidCreds
	}
	if err != nil {
		return TokenPair{}, err
	}

	if !u.IsActive {
		return TokenPair{}, ErrAccountInactive
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return TokenPair{}, ErrInvalidCreds
	}

	return s.generateTokenPair(u.ID, u.Role)
}

func (s *service) RefreshTokens(refreshToken string) (TokenPair, error) {
	if refreshToken == "" {
		return TokenPair{}, ErrTokenInvalid
	}

	tokenHash := hashToken(refreshToken)

	userID, err := s.repo.FindRefreshToken(tokenHash)
	if err != nil {
		return TokenPair{}, ErrTokenInvalid
	}

	// Revoke the old token before issuing a new pair (rotation).
	if err = s.repo.RevokeRefreshToken(tokenHash); err != nil {
		return TokenPair{}, err
	}

	u, err := s.repo.FindByID(userID)
	if err != nil {
		return TokenPair{}, err
	}

	return s.generateTokenPair(u.ID, u.Role)
}

func (s *service) Logout(refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.repo.RevokeRefreshToken(hashToken(refreshToken))
}

func (s *service) Me(userID string) (User, error) {
	return s.repo.FindByID(userID)
}

func (s *service) UpdateProfile(userID string, req UpdateProfileRequest) (User, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return User{}, err
	}

	name := strings.TrimSpace(req.Name)
	if name != "" {
		if len(name) < 2 || len(name) > 100 {
			return User{}, errors.New("name must be between 2 and 100 characters")
		}
		if err = s.repo.Update(userID, name); err != nil {
			return User{}, err
		}
		u.Name = name
	}

	if req.NewPassword != "" {
		if req.OldPassword == "" {
			return User{}, errors.New("old_password is required to change password")
		}
		if len(req.NewPassword) < 8 {
			return User{}, errors.New("new password must be at least 8 characters")
		}
		_, currentHash, err := s.repo.FindByEmail(u.Email)
		if err != nil {
			return User{}, err
		}
		if err = bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword)); err != nil {
			return User{}, ErrWrongPassword
		}
		newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
		if err != nil {
			return User{}, err
		}
		if err = s.repo.UpdatePassword(userID, string(newHash)); err != nil {
			return User{}, err
		}
	}

	return u, nil
}

func (s *service) ListUsers(page, limit int, role string, active *bool) ([]User, int, error) {
	return s.repo.ListUsers(page, limit, role, active)
}

func (s *service) SetActive(userID, callerID string, active bool) error {
	if userID == callerID {
		return ErrSelfModify
	}
	return s.repo.SetActive(userID, active)
}

func (s *service) Promote(userID string) error {
	return s.repo.SetRole(userID, "admin")
}

// ─────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────

func (s *service) generateTokenPair(userID, role string) (TokenPair, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(s.jwtAccessTTL).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return TokenPair{}, err
	}

	rawRefresh := uuid.New().String()
	tokenHash := hashToken(rawRefresh)
	expiresAt := now.Add(s.jwtRefreshTTL)

	if err = s.repo.SaveRefreshToken(userID, tokenHash, expiresAt); err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessStr,
		RefreshToken: rawRefresh,
	}, nil
}

// hashToken returns the hex-encoded SHA-256 hash of a token string.
// Refresh tokens are stored hashed so a DB breach cannot replay them.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

//go:build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/iqbaleff214/kamus-banjar-api/internal/config"
	"github.com/iqbaleff214/kamus-banjar-api/internal/seeder"
	"github.com/iqbaleff214/kamus-banjar-api/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testApp spins up the full Fiber app backed by a real test DB.
// Requires TEST_MYSQL_DSN env var.
func testApp(t *testing.T) *httptest.Server {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	require.NotEmpty(t, dsn, "TEST_MYSQL_DSN must be set for integration tests")

	cfg := config.Config{
		Port:          ":0",
		MySQLDSN:      dsn,
		JWTSecret:     "integration-test-secret-key-32chars!!",
		JWTAccessTTL:  15 * time.Minute,
		JWTRefreshTTL: 168 * time.Hour,
	}

	db := cfg.OpenDB()
	t.Cleanup(func() { db.Close() })

	seeder.Seed(db, nil) // no-op: test DB seeded separately

	app := server.New(db, cfg)
	srv := httptest.NewServer(app)
	t.Cleanup(srv.Close)
	return srv
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

func post(t *testing.T, srv *httptest.Server, path string, body any, token string) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func get(t *testing.T, srv *httptest.Server, path, token string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, srv.URL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decode(t *testing.T, r *http.Response, v any) {
	t.Helper()
	defer r.Body.Close()
	b, _ := io.ReadAll(r.Body)
	require.NoError(t, json.Unmarshal(b, v))
}

func randomEmail() string {
	return fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
}

// ─────────────────────────────────────────────────────────────
// Tests
// ─────────────────────────────────────────────────────────────

// TestFlow_RegisterLoginGetToken verifies register → login → get profile.
func TestFlow_RegisterLoginGetToken(t *testing.T) {
	srv := testApp(t)
	email := randomEmail()

	// 1. Register
	regResp := post(t, srv, "/api/v1/auth/register", map[string]any{
		"name":     "Integration User",
		"email":    email,
		"password": "password1234",
	}, "")
	assert.Equal(t, http.StatusCreated, regResp.StatusCode)

	// 2. Login
	loginResp := post(t, srv, "/api/v1/auth/login", map[string]any{
		"email":    email,
		"password": "password1234",
	}, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	var loginBody struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	decode(t, loginResp, &loginBody)
	require.NotEmpty(t, loginBody.Data.AccessToken)

	// 3. GET /auth/me
	meResp := get(t, srv, "/api/v1/auth/me", loginBody.Data.AccessToken)
	assert.Equal(t, http.StatusOK, meResp.StatusCode)
}

// TestFlow_SubmitWordContribution verifies register → login → submit word → verify pending.
func TestFlow_SubmitWordContribution(t *testing.T) {
	srv := testApp(t)
	email := randomEmail()

	post(t, srv, "/api/v1/auth/register", map[string]any{
		"name": "Contributor", "email": email, "password": "password1234",
	}, "")

	loginResp := post(t, srv, "/api/v1/auth/login", map[string]any{
		"email": email, "password": "password1234",
	}, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	var loginBody struct {
		Data struct{ AccessToken string `json:"access_token"` } `json:"data"`
	}
	decode(t, loginResp, &loginBody)
	token := loginBody.Data.AccessToken

	// Submit a new word
	submitResp := post(t, srv, "/api/v1/contributions", map[string]any{
		"word":     fmt.Sprintf("testword_%d", time.Now().UnixNano()),
		"alphabet": "t",
		"meanings": []map[string]any{
			{"definitions": []map[string]any{
				{"definition": "kata uji", "partOfSpeech": "n"},
			}},
		},
		"derivatives": []any{},
	}, token)
	assert.Equal(t, http.StatusCreated, submitResp.StatusCode)

	// Verify it appears in /contributions/mine
	mineResp := get(t, srv, "/api/v1/contributions/mine", token)
	assert.Equal(t, http.StatusOK, mineResp.StatusCode)
}

// TestFlow_InvalidLogin verifies wrong password is rejected.
func TestFlow_InvalidLogin(t *testing.T) {
	srv := testApp(t)
	email := randomEmail()

	post(t, srv, "/api/v1/auth/register", map[string]any{
		"name": "Bad Login", "email": email, "password": "password1234",
	}, "")

	resp := post(t, srv, "/api/v1/auth/login", map[string]any{
		"email": email, "password": "wrongpassword",
	}, "")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestFlow_UnauthenticatedContributionRejected verifies that unauthenticated
// contribution submission returns 401.
func TestFlow_UnauthenticatedContributionRejected(t *testing.T) {
	srv := testApp(t)

	resp := post(t, srv, "/api/v1/contributions", map[string]any{
		"word": "unauthorised", "alphabet": "u",
	}, "")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestFlow_PublicDictionaryEndpoints verifies public endpoints work without auth.
func TestFlow_PublicDictionaryEndpoints(t *testing.T) {
	srv := testApp(t)

	alphabetsResp := get(t, srv, "/api/v1/alphabets", "")
	assert.Equal(t, http.StatusOK, alphabetsResp.StatusCode)
}

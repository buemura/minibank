package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestManager() *Manager {
	return NewManager("test-secret-key", 15*time.Minute, 7*24*time.Hour)
}

func TestGenerateTokenPair_Success(t *testing.T) {
	m := newTestManager()

	pair, err := m.GenerateTokenPair("user-123", "test@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, int64(900), pair.ExpiresIn) // 15 minutes = 900 seconds
}

func TestValidateToken_ValidToken(t *testing.T) {
	m := newTestManager()

	pair, err := m.GenerateTokenPair("user-123", "test@example.com")
	require.NoError(t, err)

	claims, err := m.ValidateToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "minibank", claims.Issuer)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	m := NewManager("test-secret-key", -1*time.Hour, 7*24*time.Hour)

	pair, err := m.GenerateTokenPair("user-123", "test@example.com")
	require.NoError(t, err)

	_, err = m.ValidateToken(pair.AccessToken)
	assert.ErrorIs(t, err, ErrExpiredToken)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	m1 := NewManager("secret-key-1", 15*time.Minute, 7*24*time.Hour)
	m2 := NewManager("secret-key-2", 15*time.Minute, 7*24*time.Hour)

	pair, err := m1.GenerateTokenPair("user-123", "test@example.com")
	require.NoError(t, err)

	_, err = m2.ValidateToken(pair.AccessToken)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestValidateToken_MalformedToken(t *testing.T) {
	m := newTestManager()

	_, err := m.ValidateToken("not-a-valid-token")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestGenerateAccessToken_Success(t *testing.T) {
	m := newTestManager()

	token, err := m.GenerateAccessToken("user-456", "user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := m.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "user-456", claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	m := newTestManager()

	_, err := m.ValidateToken("")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

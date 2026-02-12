package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash_Success(t *testing.T) {
	hash, err := Hash("mypassword123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "mypassword123", hash)
}

func TestVerify_CorrectPassword(t *testing.T) {
	hash, err := Hash("mypassword123")
	require.NoError(t, err)

	assert.True(t, Verify("mypassword123", hash))
}

func TestVerify_WrongPassword(t *testing.T) {
	hash, err := Hash("mypassword123")
	require.NoError(t, err)

	assert.False(t, Verify("wrongpassword", hash))
}

func TestHash_DifferentHashesForSameInput(t *testing.T) {
	hash1, err := Hash("samepassword")
	require.NoError(t, err)

	hash2, err := Hash("samepassword")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "bcrypt should produce different hashes due to salt")
}

func TestHash_EmptyPassword(t *testing.T) {
	hash, err := Hash("")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, Verify("", hash))
}

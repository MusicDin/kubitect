package keygen

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyWrite(t *testing.T) {
	keyPath := path.Join(t.TempDir(), "id_rsa")
	key := key{
		value: []byte("test"),
	}

	err := key.Write(keyPath)
	require.NoError(t, err)

	keyFile, err := os.ReadFile(keyPath)
	require.NoError(t, err)
	assert.Equal(t, "test", string(keyFile))
}

func TestKeyWrite_InvalidPath(t *testing.T) {
	keyPath := path.Join(t.TempDir(), "invalid", "id_rsa")
	key := key{
		value: []byte("test"),
	}

	err := key.Write(keyPath)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestGeneratePrivateKey(t *testing.T) {
	key, err := generatePrivateKey(1024)
	require.NoError(t, err)
	assert.NotEmpty(t, key)
}

func TestEncodePrivateKey(t *testing.T) {
	key, err := generatePrivateKey(1024)
	require.NoError(t, err)

	pem := encodePrivateKey(key)
	assert.Contains(t, string(pem), "RSA PRIVATE KEY")
}

func TestGeneratePublicKey(t *testing.T) {
	key, err := generatePrivateKey(1024)
	require.NoError(t, err)

	publicKey, err := generatePublicKey(key)
	require.NoError(t, err)
	assert.NotEmpty(t, publicKey)
}

func TestGeneratePublicKey_NilPrivateKey(t *testing.T) {
	_, err := generatePublicKey(nil)
	assert.EqualError(t, err, "generate public key: received <nil> private key")
}

func TestNewKeyPair(t *testing.T) {
	_, err := NewKeyPair(1024)
	assert.NoError(t, err)
}

func TestReadKeyPair(t *testing.T) {
	tmpDir := t.TempDir()

	kp1, err := NewKeyPair(1024)
	require.NoError(t, err)
	require.NoError(t, kp1.Write(tmpDir, "key"))

	// Existing key should be fetched
	kp2, err := ReadKeyPair(tmpDir, "key")
	require.NoError(t, err)
	assert.Equal(t, kp1, kp2)
}

func TestReadKeyPair_FailReadingPrivateKey(t *testing.T) {
	tmpDir := t.TempDir()
	keyName := "key"
	keyPath := path.Join(tmpDir, keyName)

	kp, err := NewKeyPair(1024)
	require.NoError(t, err)
	require.NoError(t, kp.Write(tmpDir, keyName))

	// Make an empty directory on a public key path
	require.NoError(t, os.Remove(keyPath))
	require.NoError(t, os.Mkdir(keyPath, 0700))

	_, err = ReadKeyPair(tmpDir, keyName)
	assert.ErrorContains(t, err, NewKeyFileError("private", keyName, errors.New("")).Error())
}

func TestReadKeyPair_FailReadingPublicKey(t *testing.T) {
	tmpDir := t.TempDir()
	keyName := "key"
	keyPath := path.Join(tmpDir, keyName+".pub")

	kp, err := NewKeyPair(1024)
	require.NoError(t, err)
	require.NoError(t, kp.Write(tmpDir, keyName))

	// Make an empty directory on a public key path
	require.NoError(t, os.Remove(keyPath))
	require.NoError(t, os.Mkdir(keyPath, 0700))

	_, err = ReadKeyPair(tmpDir, keyName)
	assert.ErrorContains(t, err, NewKeyFileError("public", keyName+".pub", errors.New("")).Error())
}

func TestNewKeyPair_Invalid(t *testing.T) {
	_, err := NewKeyPair(-1)
	assert.ErrorContains(t, err, "generate private key: crypto")
}

func TestKeyPair_Write(t *testing.T) {
	tmpDir := t.TempDir()
	keyName := "key"

	kp, err := NewKeyPair(1024)
	require.NoError(t, err)
	require.NoError(t, kp.Write(tmpDir, keyName))

	privKeyPath := path.Join(tmpDir, keyName)
	privKey, err := os.ReadFile(privKeyPath)
	require.NoError(t, err)
	require.Contains(t, string(privKey), "RSA PRIVATE KEY")

	pubKeyPath := path.Join(tmpDir, keyName+".pub")
	pubKey, err := os.ReadFile(pubKeyPath)
	require.NoError(t, err)
	assert.Contains(t, string(pubKey), "ssh-rsa")
}

func TestKeyPair_Write_InvalidPath(t *testing.T) {
	kp, err := NewKeyPair(1024)
	require.NoError(t, err)

	err = kp.Write(t.TempDir(), path.Join("invalid", "key"))
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestKeyPairExists(t *testing.T) {
	keyDir := t.TempDir()
	keyName := "key"
	kp, err := NewKeyPair(1024)
	require.NoError(t, err)
	require.NoError(t, kp.Write(keyDir, keyName))
	assert.True(t, KeyPairExists(keyDir, keyName), "KeyPair does not exists after being generated")
}

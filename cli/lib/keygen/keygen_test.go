package keygen

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyWrite(t *testing.T) {
	keyPath := path.Join(t.TempDir(), "id_rsa")
	key := key{
		value: []byte("test"),
	}

	err := key.Write(keyPath)
	assert.NoError(t, err)

	keyFile, err := ioutil.ReadFile(keyPath)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(keyFile))
}

func TestGeneratePrivateKey(t *testing.T) {
	key, err := generatePrivateKey(512)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
}

func TestGeneratePrivateKey_Error(t *testing.T) {
	_, err := generatePrivateKey(-1)
	assert.EqualError(t, err, "generate private key: crypto/rsa: too few primes of given length to generate an RSA key")
}

func TestEncodePrivateKey(t *testing.T) {
	key, err := generatePrivateKey(512)
	assert.NoError(t, err)

	pem := encodePrivateKey(key)
	assert.NotEmpty(t, pem)
	assert.Contains(t, string(pem), "RSA PRIVATE KEY")
}

func TestGeneratePublicKey(t *testing.T) {
	key, err := generatePrivateKey(512)
	assert.NoError(t, err)

	publicKey, err := generatePublicKey(key)
	assert.NoError(t, err)
	assert.NotEmpty(t, publicKey)
}

func TestGeneratePublicKey_NilPrivateKey(t *testing.T) {
	_, err := generatePublicKey(nil)
	assert.EqualError(t, err, "generate public key: received <nil> private key")
}

func TestNewKeyPair(t *testing.T) {
	_, err := NewKeyPair(512)
	assert.NoError(t, err)
}

func TestNewKeyPair_Invalid(t *testing.T) {
	_, err := NewKeyPair(-1)
	assert.ErrorContains(t, err, "generate private key: crypto")
}

func TestWriteKeyPair(t *testing.T) {
	kpPath := t.TempDir()
	kp, err := NewKeyPair(512)
	assert.NoError(t, err)

	err = kp.WriteKeys(kpPath)
	assert.NoError(t, err)

	privKeyPath := path.Join(kpPath, "id_rsa")
	privKey, err := ioutil.ReadFile(privKeyPath)
	assert.NoError(t, err)
	assert.Contains(t, string(privKey), "RSA PRIVATE KEY")

	pubKeyPath := path.Join(kpPath, "id_rsa.pub")
	pubKey, err := ioutil.ReadFile(pubKeyPath)
	assert.NoError(t, err)
	assert.Contains(t, string(pubKey), "ssh-rsa")
}

func TestWriteKeyPair_InvalidPath(t *testing.T) {
	kp, err := NewKeyPair(512)
	assert.NoError(t, err)

	err = kp.WriteKeys("#")
	assert.EqualError(t, err, "open #/id_rsa: no such file or directory")
}

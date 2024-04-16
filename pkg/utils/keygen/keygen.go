package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

type (
	Key interface {
		Write(path string) error
	}

	key struct {
		value []byte
	}
)

// Write writes the key to the specified path.
// Note that the directory to which the key is
// written must exist.
func (k key) Write(path string) error {
	return os.WriteFile(path, k.value, 0600)
}

type (
	KeyPair interface {
		PublicKey() Key
		PrivateKey() Key
		Write(dir, keyName string) error
	}

	keyPair struct {
		private key
		public  key
	}
)

func (p keyPair) PrivateKey() Key {
	return p.private
}

func (p keyPair) PublicKey() Key {
	return p.public
}

// Write keys creates the keys directory and writes the keys to it.
func (p keyPair) Write(dir, keyName string) error {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	privKeyPath := path.Join(dir, keyName)
	pubKeyPath := path.Join(dir, keyName+".pub")

	err = p.PrivateKey().Write(privKeyPath)
	if err != nil {
		return err
	}

	return p.PublicKey().Write(pubKeyPath)
}

// NewKeyPair generates new private and public key pair.
func NewKeyPair(bitSize int) (KeyPair, error) {
	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return nil, err
	}

	pair := keyPair{}
	pair.private.value = encodePrivateKey(privateKey)
	pair.public.value, err = generatePublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	return pair, nil
}

// ReadKeyPair reads public and private keys with a given name
// from the specified directory.
func ReadKeyPair(dir, keyName string) (KeyPair, error) {
	var err error
	var pair keyPair

	privKeyPath := path.Join(dir, keyName)
	pair.private.value, err = os.ReadFile(privKeyPath)
	if err != nil {
		return nil, NewKeyFileError("private", keyName, err)
	}

	pubKeyName := keyName + ".pub"
	pubKeyPath := path.Join(dir, pubKeyName)
	pair.public.value, err = os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, NewKeyFileError("public", keyName+".pub", err)
	}

	return pair, nil
}

// Returns true if key pair with a given name exists in a
// specified directory.
func KeyPairExists(dir, keyName string) bool {
	kp, err := ReadKeyPair(dir, keyName)
	return err == nil && kp != nil
}

// generatePrivateKey creates a RSA private key of specified bit size.
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	pk, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, fmt.Errorf("generate private key: %v", err)
	}

	return pk, pk.Validate()
}

// encodePrivateKey encodes Private Key from RSA to PEM format.
func encodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	pemBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.EncodeToMemory(&pemBlock)
}

// generatePublicKey generates a new public key from a public part of the
// private key. The returned bytes are suitable for writing a .pub file,
// since they are in format "ssh-rsa ...".
func generatePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("generate public key: received <nil> private key")
	}

	publicRsaKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(publicRsaKey), nil
}

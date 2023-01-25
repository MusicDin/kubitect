package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
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

func (k key) Write(path string) error {
	return ioutil.WriteFile(path, k.value, 0600)
}

type (
	KeyPair interface {
		WriteKeys(dir string) error
		PublicKey() Key
		PrivateKey() Key
	}

	keyPair struct {
		private Key
		public  Key
	}
)

func (p keyPair) PrivateKey() Key {
	return p.private
}

func (p keyPair) PublicKey() Key {
	return p.public
}

func NewKeyPair(bitSize int) (KeyPair, error) {
	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := encodePrivateKey(privateKey)
	publicKeyBytes, err := generatePublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	pair := keyPair{
		private: key{
			value: privateKeyBytes,
		},
		public: key{
			value: publicKeyBytes,
		},
	}

	return pair, nil
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

func (p keyPair) WriteKeys(dir string) error {
	err := p.PrivateKey().Write(path.Join(dir, "id_rsa"))
	if err != nil {
		return err
	}

	return p.PublicKey().Write(path.Join(dir, "id_rsa.pub"))
}

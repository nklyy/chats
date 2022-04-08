package keyPair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type KeyPair interface {
	GenerateKeyPair() ([]byte, []byte, error)
	WriteKeysToKeysFolder(roomName string, privKey, pubKey []byte) error
}

type keyPair struct {
	logger *zap.SugaredLogger
}

func NewKeyPairService(logger *zap.SugaredLogger) (KeyPair, error) {
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &keyPair{logger: logger}, nil
}

func (k *keyPair) GenerateKeyPair() ([]byte, []byte, error) {
	bitSize := 4096

	// Generate RSA key
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		k.logger.Errorf("failed to generate keys %v", err)
		return nil, nil, err
	}

	// Extract pub key
	pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM
	privPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	return privPEM, pubPEM, nil
}

func (k *keyPair) WriteKeysToKeysFolder(roomName string, privKey, pubKey []byte) error {
	err := os.Mkdir("keys/"+roomName, 0755)
	if err != nil {
		k.logger.Errorf("failed to create roomKeys folder %v", err)
		return err
	}

	// Write private key
	privPath := fmt.Sprintf("keys/%v/%v.rsa", roomName, roomName)
	if err := ioutil.WriteFile(privPath, privKey, 0700); err != nil {
		k.logger.Errorf("failed to save private key %v", err)
		return err
	}

	// Write private key
	pubPath := fmt.Sprintf("keys/%v/%v.rsa.pub", roomName, roomName)
	if err := ioutil.WriteFile(pubPath, pubKey, 0700); err != nil {
		k.logger.Errorf("failed to save public key %v", err)
		return err
	}

	return nil
}

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 16
	keySize    = 32 // AES-256
	iterations = 100000
)

type CryptoManager struct {
	key []byte
}

// NewCryptoManager crea un gestor con una clave derivada de una contraseña
func NewCryptoManager(password string, salt []byte) *CryptoManager {
	if salt == nil {
		salt = []byte("anongo-default-salt") // Se recomienda generar uno aleatorio y guardarlo
	}
	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
	return &CryptoManager{key: key}
}

// Encrypt cifra datos usando AES-GCM
func (cm *CryptoManager) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(cm.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Sella los datos: nonce + ciphertext + tag
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt descifra datos cifrados con Encrypt
func (cm *CryptoManager) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(cm.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext demasiado corto")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, actualCiphertext, nil)
}

package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	password := "MiSuperPasswordSegura123!"
	plaintext := []byte("Mensaje secreto de AnonGo")
	
	cm := NewCryptoManager(password, nil)

	// 1. Probar Cifrado
	ciphertext, err := cm.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Fallo al cifrar: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("El ciphertext es igual al plaintext (no hubo cifrado real)")
	}

	// 2. Probar Descifrado
	decrypted, err := cm.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Fallo al descifrar: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Descifrado incorrecto. Esperaba %s, obtuve %s", plaintext, decrypted)
	}
}

func TestWrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"
	plaintext := []byte("Información sensible")

	cmCorrect := NewCryptoManager(password, nil)
	cmWrong := NewCryptoManager(wrongPassword, nil)

	ciphertext, _ := cmCorrect.Encrypt(plaintext)

	// Intentar descifrar con clave incorrecta (AES-GCM debería fallar en la autenticación)
	_, err := cmWrong.Decrypt(ciphertext)
	if err == nil {
		t.Error("El descifrado tuvo éxito con una contraseña incorrecta; esto es un fallo de seguridad grave.")
	}
}

func TestNonceUniqueness(t *testing.T) {
	password := "password"
	plaintext := []byte("Mismo mensaje")
	cm := NewCryptoManager(password, nil)

	c1, _ := cm.Encrypt(plaintext)
	c2, _ := cm.Encrypt(plaintext)

	if bytes.Equal(c1, c2) {
		t.Error("El cifrado generó el mismo resultado para el mismo mensaje (el nonce no es aleatorio/único)")
	}
}

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
)

func loadECDSAKey(filePath string) (*ecdsa.PrivateKey, error) {
	keyJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var keyData ECDSAKey
	if err := json.Unmarshal(keyJSON, &keyData); err != nil {
		return nil, err
	}

	d, _ := new(big.Int).SetString(keyData.D, 16)
	x, _ := new(big.Int).SetString(keyData.X, 16)
	y, _ := new(big.Int).SetString(keyData.Y, 16)

	privateKey := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		},
	}

	return privateKey, nil
}

func loadRSAKey(filePath string) (*rsa.PrivateKey, error) {
	keyJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var keyData RSAKey
	if err := json.Unmarshal(keyJSON, &keyData); err != nil {
		return nil, err
	}

	n := new(big.Int)
	n.SetString(keyData.N, 16)

	d := new(big.Int)
	d.SetString(keyData.D, 16)

	privateKey := &rsa.PrivateKey{
		D: d,
		PublicKey: rsa.PublicKey{
			N: n,
			E: 65537,
		},
	}

	return privateKey, nil
}

func main() {
	ecdsaKeyPath := "ecdsa_key.json"
	rsaKeyPath := "rsa_key.json"
	filePath := flag.String("file", "", "Путь к файлу")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Укажите файл -file")
		return
	}

	plaintext, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	rsaPrivateKey, err := loadRSAKey(rsaKeyPath)
	if err != nil {
		fmt.Println("Ошибка загруки RSA ключа, генерация нового:", err)
		rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Println("Ошибка генерации RSA:", err)
			return
		}
	}

	ecdsaPrivateKey, err := loadECDSAKey(ecdsaKeyPath)
	if err != nil {
		fmt.Println("Ошибка загрузки ECDSA ключа, генерация нового:", err)
		ecdsaPrivateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			fmt.Println("Ошибка генерации ECDSA ключа:", err)
			return
		}
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &rsaPrivateKey.PublicKey, plaintext, nil)
	if err != nil {
		fmt.Println("Ошибка шифровки:", err)
		return
	}
	fmt.Println("Зашифрованный текст:", ciphertext)

	hash := sha256.Sum256(ciphertext)
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivateKey, hash[:])
	if err != nil {
		fmt.Println("Ошибка подписи:", err)
		return
	}
	fmt.Println("Подпись:", r, s)

	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, ciphertext, nil)
	if err != nil {
		fmt.Println("Ошибка дешифровки:", err)
		return
	}
	fmt.Println("Дешифрованный текст:", string(decryptedText))

}

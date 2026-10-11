package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func CalcHash(key, value string) (string, error) {

	if key == "" {
		return "", fmt.Errorf("ключ не может быть пустым")
	}
	if value == "" {
		return "", fmt.Errorf("значение не может быть пустым")
	}

	data := []byte(value + key)
	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:]), nil
}

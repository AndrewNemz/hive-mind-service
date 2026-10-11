package handlersutils

import (
	"fmt"
	"hiv_mind/pkg/hash"
)

func IsEqualHash(requestHash, hashK string, body []byte) error {

	if hashK == "" {
		return nil
	}
	if requestHash == "" {
		return fmt.Errorf("Отсутвтует значения заголовка HashSHA256!")
	}

	serverHash, err := hash.CalcHash(hashK, string(body))
	if err != nil {
		return fmt.Errorf("ошибка вычисления хеша на сервере: %w", err)
	}

	if serverHash != requestHash {
		return fmt.Errorf("Ошибка! Хэши не совпадают!")
	}
	return nil
}

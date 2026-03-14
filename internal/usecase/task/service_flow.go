package task

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func defaultIdentifierFunction() (string, error) {
	randomBytes := make([]byte, 16)
	_, randomError := rand.Read(randomBytes)
	if randomError != nil {
		return "", fmt.Errorf("identifier random read failure: %w", randomError)
	}
	return hex.EncodeToString(randomBytes), nil
}

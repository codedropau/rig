package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
)

// Base64 encode the Docker registry authentication credentials.
func Base64(username, password string) (string, error) {
	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}

	authBytes, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth: %w", err)
	}

	return base64.URLEncoding.EncodeToString(authBytes), nil
}

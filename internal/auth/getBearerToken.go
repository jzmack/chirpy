package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("No Authorization header")
	}
	authHeaderSlice := strings.Split(authHeader, " ")
	if authHeaderSlice[0] != "Bearer" {
		return "", fmt.Errorf("Not a Bearer")
	}
	if len(authHeaderSlice) < 2 {
		return "", fmt.Errorf("Header does not contain Token")
	}
	return strings.TrimSpace(authHeaderSlice[1]), nil
}

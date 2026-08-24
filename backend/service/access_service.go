package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var allowedColors = []string{"#E53935", "#1E88E5", "#43A047", "#FDD835", "#8E24AA", "#FB8C00"}

func RandomBackgroundColor() string {
	return allowedColors[time.Now().UnixNano()%int64(len(allowedColors))]
}

func IsAllowedBackgroundColor(color string) bool {
	for _, item := range allowedColors {
		if color == item {
			return true
		}
	}
	return false
}

// GenerateQRToken creates a signed, short-lived code. The mobile app renders this
// string as a QR image; no QR image needs to be stored on the server.
func GenerateQRToken(secret, sessionID, employeeID string, expiresAt time.Time) string {
	payload := fmt.Sprintf("%s|%s|%d", sessionID, employeeID, expiresAt.Unix())
	signature := sign(secret, payload)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + signature
}

func VerifyQRToken(secret, token string, now time.Time) (sessionID string, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", errors.New("formato de QR inválido")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("código QR inválido")
	}
	payload := string(payloadBytes)
	if subtle.ConstantTimeCompare([]byte(sign(secret, payload)), []byte(parts[1])) != 1 {
		return "", errors.New("firma de QR inválida")
	}
	fields := strings.Split(payload, "|")
	if len(fields) != 3 {
		return "", errors.New("contenido de QR inválido")
	}
	expiresUnix, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil || now.After(time.Unix(expiresUnix, 0)) {
		return "", errors.New("QR expirado")
	}
	return fields[0], nil
}

func sign(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

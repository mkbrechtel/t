package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	uuidv7 "github.com/gofrs/uuid/v5"
)

func NewUUID() (uid uuidv7.UUID) {
	return uuidv7.Must(uuidv7.NewV7())
}

func ShortEncodeUUID(uid uuidv7.UUID) string {
	encoded := base64.RawURLEncoding.EncodeToString(uid.Bytes())
	return strings.NewReplacer(
		"-", "98",
		"_", "90",
		"9", "99",
	).Replace(encoded)
}

func LongEncodeUUID(uid uuidv7.UUID) string {
	return uid.String()
}

func EncodeUUID(uid uuidv7.UUID) string {
	return ShortEncodeUUID(uid)
}

func DecodeUUID(encid string) (uuidv7.UUID, error) {
	if len(encid) == 0 {
		return uuidv7.Nil, fmt.Errorf("empty UUID string")
	}

	if len(encid) == 36 {
		u, err := uuidv7.FromString(encid)
		if err != nil {
			return uuidv7.Nil, fmt.Errorf("failed to parse UUID: %w", err)
		}
		return u, nil
	}

	// Reverse the character replacements
	decoded := strings.NewReplacer(
		"98", "-",
		"90", "_",
		"99", "9",
	).Replace(encid)

	b, err := base64.RawURLEncoding.DecodeString(decoded)
	if err != nil {
		return uuidv7.Nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	u, err := uuidv7.FromBytes(b)
	if err != nil {
		return uuidv7.Nil, fmt.Errorf("failed to parse UUID from bytes: %w", err)
	}
	return u, nil
}

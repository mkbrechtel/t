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

// Base64 alphabet shifted to start with 't' after base64url padding is removed
// Standard:  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
// Modified:  tuvwxyz0123456789-_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrs
var customEncoding = base64.NewEncoding(
	"tuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqr+-#",
).WithPadding(base64.NoPadding)

// Global replacers to avoid recreating them on every encode/decode
var (
	encodeReplacer = strings.NewReplacer( //qrs
		"+", "sp",
		"-", "sq",
		"#", "sr",
	)
	decodeReplacer = strings.NewReplacer(
		"sp", "+", 
		"sq", "-", 
		"sr", "#", 
	)
)

func ShortEncodeUUID(uid uuidv7.UUID) string {
	encoded := customEncoding.EncodeToString(uid.Bytes())
	return encodeReplacer.Replace(encoded)
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

	decoded := decodeReplacer.Replace(encid)

	b, err := customEncoding.DecodeString(decoded)
	if err != nil {
		return uuidv7.Nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	u, err := uuidv7.FromBytes(b)
	if err != nil {
		return uuidv7.Nil, fmt.Errorf("failed to parse UUID from bytes: %w", err)
	}
	return u, nil
}

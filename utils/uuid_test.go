package utils

import (
	"testing"

	uuidv7 "github.com/gofrs/uuid/v5"
)

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	if uuid == uuidv7.Nil {
		t.Error("NewUUID() returned nil UUID")
	}
}

func TestUUIDEncodingDecoding(t *testing.T) {
	// Setup test UUID
	staticUUID := uuidv7.Must(uuidv7.FromString("0192d9f7-edc6-76d8-8aae-c9b5b0237c0d"))

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Test Short Encoding and Decoding",
			testFunc: func(t *testing.T) {
				encoded := ShortEncodeUUID(staticUUID)
				decoded, err := DecodeUUID(encoded)
				if err != nil {
					t.Errorf("Failed to decode short UUID: %v", err)
				}
				if decoded != staticUUID {
					t.Errorf("Short encode/decode mismatch: got %v, want %v", decoded, staticUUID)
				}
			},
		},
		{
			name: "Test Long Encoding and Decoding",
			testFunc: func(t *testing.T) {
				encoded := LongEncodeUUID(staticUUID)
				decoded, err := DecodeUUID(encoded)
				if err != nil {
					t.Errorf("Failed to decode long UUID: %v", err)
				}
				if decoded != staticUUID {
					t.Errorf("Long encode/decode mismatch: got %v, want %v", decoded, staticUUID)
				}
			},
		},
		{
			name: "Test Default Encoding",
			testFunc: func(t *testing.T) {
				encoded := EncodeUUID(staticUUID)
				decoded, err := DecodeUUID(encoded)
				if err != nil {
					t.Errorf("Failed to decode default encoded UUID: %v", err)
				}
				if decoded != staticUUID {
					t.Errorf("Default encode/decode mismatch: got %v, want %v", decoded, staticUUID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestDecodeUUIDErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Invalid base64",
			input:   "!invalid-base64!",
			wantErr: true,
		},
		{
			name:    "Invalid UUID format",
			input:   "not-a-valid-uuid-format",
			wantErr: true,
		},
		{
			name:    "Invalid UUID string length",
			input:   "0192d9f7-edc6-76d8-8aae-c9b5b0237c0", // Missing last character
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeUUID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUUIDRoundTrip(t *testing.T) {
	// Generate a new UUID and test full round-trip
	original := NewUUID()
	
	// Test short encoding round-trip
	shortEncoded := ShortEncodeUUID(original)
	shortDecoded, err := DecodeUUID(shortEncoded)
	if err != nil {
		t.Errorf("Failed to decode short encoded UUID: %v", err)
	}
	if shortDecoded != original {
		t.Errorf("Short encoding round-trip failed: got %v, want %v", shortDecoded, original)
	}

	// Test long encoding round-trip
	longEncoded := LongEncodeUUID(original)
	longDecoded, err := DecodeUUID(longEncoded)
	if err != nil {
		t.Errorf("Failed to decode long encoded UUID: %v", err)
	}
	if longDecoded != original {
		t.Errorf("Long encoding round-trip failed: got %v, want %v", longDecoded, original)
	}
}

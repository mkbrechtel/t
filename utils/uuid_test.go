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

func TestStaticUUIDTransformations(t *testing.T) {
	// Static test values
	expectedUUID := "0192d9fd-d725-765d-8305-a5cdcd3ffe63"
	expectedShortEncoding := "AZLZ_dcldl2DBaXNzT_-Yw"
	expectedLongEncoding := expectedUUID

	// Create our base UUID
	baseUUID := uuidv7.Must(uuidv7.FromString(expectedUUID))

	// Test cases for all possible transformations
	t.Run("UUID to Short Encoding", func(t *testing.T) {
		shortEncoded := ShortEncodeUUID(baseUUID)
		if shortEncoded != expectedShortEncoding {
			t.Errorf("ShortEncodeUUID() = %v, want %v", shortEncoded, expectedShortEncoding)
		}
	})

	t.Run("UUID to Long Encoding", func(t *testing.T) {
		longEncoded := LongEncodeUUID(baseUUID)
		if longEncoded != expectedLongEncoding {
			t.Errorf("LongEncodeUUID() = %v, want %v", longEncoded, expectedLongEncoding)
		}
	})

	t.Run("Short Encoding to UUID", func(t *testing.T) {
		decoded, err := DecodeUUID(expectedShortEncoding)
		if err != nil {
			t.Errorf("DecodeUUID(short) failed: %v", err)
		}
		if decoded != baseUUID {
			t.Errorf("DecodeUUID(short) = %v, want %v", decoded, baseUUID)
		}
	})

	t.Run("Long Encoding to UUID", func(t *testing.T) {
		decoded, err := DecodeUUID(expectedLongEncoding)
		if err != nil {
			t.Errorf("DecodeUUID(long) failed: %v", err)
		}
		if decoded != baseUUID {
			t.Errorf("DecodeUUID(long) = %v, want %v", decoded, baseUUID)
		}
	})

	t.Run("Default Encode matches Short Encode", func(t *testing.T) {
		encoded := EncodeUUID(baseUUID)
		if encoded != expectedShortEncoding {
			t.Errorf("EncodeUUID() = %v, want %v", encoded, expectedShortEncoding)
		}
	})

	// Test round-trip transformations
	t.Run("Short Encoding Round Trip", func(t *testing.T) {
		// Short -> UUID -> Short
		decoded, err := DecodeUUID(expectedShortEncoding)
		if err != nil {
			t.Errorf("DecodeUUID(short) failed: %v", err)
		}
		reencoded := ShortEncodeUUID(decoded)
		if reencoded != expectedShortEncoding {
			t.Errorf("Short round-trip failed: got %v, want %v", reencoded, expectedShortEncoding)
		}
	})

	t.Run("Long Encoding Round Trip", func(t *testing.T) {
		// Long -> UUID -> Long
		decoded, err := DecodeUUID(expectedLongEncoding)
		if err != nil {
			t.Errorf("DecodeUUID(long) failed: %v", err)
		}
		reencoded := LongEncodeUUID(decoded)
		if reencoded != expectedLongEncoding {
			t.Errorf("Long round-trip failed: got %v, want %v", reencoded, expectedLongEncoding)
		}
	})
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
			input:   "0192d9fd-d725-765d-8305-a5cdcd3ffe", // Missing last characters
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

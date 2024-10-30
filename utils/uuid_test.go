package utils

import (
	"math/rand"
	"testing"
	"time"

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
	expectedUUID := "0192da75-c158-7d7f-be3c-d5b647bf7fa8"
	expectedShortEncoding := "tI4JMLyHOGsqsq86FlAqspsrZt"
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

func TestEdgeCaseUUIDTransformations(t *testing.T) {
	cases := []struct {
		name     string
		uuid     string
		encoding string
	}{
		{
			name:     "Many special chars",
			uuid:     "0192da73-39ce-76ac-826b-bb3fd7e9fd84",
			encoding: "tI4JLiW7MZhvJqbsrksqWspQt",
		},
		{
			name:     "Sequential special chars",
			uuid:     "ffffffff-ffff-7fff-bfff-ffffffffffff",
			encoding: "srsrsrsrsrsrsrsrOsrsqsrsrsrsrsrsrsrsrsrsrf",
		},
		// {
		// 	name:     "Multiple s chars",
		// 	uuid:     "00000000-0000-0000-0000-000000000000",
		// 	encoding: "pspspspspspspspspspspspspspspspps",
		// },
		// {
		// 	name:     "Mix of special replacements",
		// 	uuid:     "12345678-90ab-7cde-f123-456789abcdef",
		// 	encoding: "psqsrspsqsrspsqsrspsqsrspsqsrsps",
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baseUUID := uuidv7.Must(uuidv7.FromString(tt.uuid))

			// Test encoding
			encoded := ShortEncodeUUID(baseUUID)
			if encoded != tt.encoding {
				t.Errorf("ShortEncodeUUID() = %v, want %v", encoded, tt.encoding)
			}

			// Test decoding
			decoded, err := DecodeUUID(tt.encoding)
			if err != nil {
				t.Errorf("DecodeUUID() failed: %v", err)
			}
			if decoded != baseUUID {
				t.Errorf("DecodeUUID() = %v, want %v", decoded, baseUUID)
			}
		})
	}
}

func TestTemporalEdgeCases(t *testing.T) {
	cases := []struct {
		name string
		time time.Time
	}{
		{
			name: "Unix epoch start",
			time: time.Unix(0, 0),
		},
		{
			name: "Current time",
			time: time.Now(),
		},
		{
			name: "Far future",
			time: time.Date(2999, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "Leap year",
			time: time.Date(2024, 2, 29, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "Millisecond precision edge",
			time: time.Date(2024, 1, 1, 0, 0, 0, 999999, time.UTC),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			uid := uuidv7.Must(uuidv7.NewV7AtTime(tt.time))
			encoded := ShortEncodeUUID(uid)
			
			decoded, err := DecodeUUID(encoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			if decoded != uid {
				t.Errorf("Round trip failed: got %v, want %v", decoded, uid)
			}

			// Verify timestamp preservation (to millisecond precision)
			originalTs, _ := uuidv7.TimestampFromV7(uid)
			decodedTs, _ := uuidv7.TimestampFromV7(decoded)
			
			originalTime, _ := originalTs.Time()
			decodedTime, _ := decodedTs.Time()

			if !originalTime.Equal(decodedTime) {
				t.Errorf("Timestamp mismatch: got %v, want %v", decodedTime, originalTime)
			}
		})
	}
}

func TestTemporalStressTest(t *testing.T) {
	const numOperations = 1_000_000

	// Start from Unix epoch
	startTime := time.Unix(0, 0)
	// Add 1 millisecond per iteration to test different timestamps
	timeIncrement := time.Millisecond

	var lastTime time.Time
	for i := 0; i < numOperations; i++ {
		testTime := startTime.Add(timeIncrement * time.Duration(i))
		
		// Create UUID with specific timestamp
		uid := uuidv7.Must(uuidv7.NewV7AtTime(testTime))
		
		// Convert to short encoding and back
		encoded := ShortEncodeUUID(uid)
		decoded, err := DecodeUUID(encoded)
		if err != nil {
			t.Fatalf("Failed to decode at iteration %d: %v", i, err)
		}

		if decoded != uid {
			t.Errorf("Round trip failed at iteration %d: got %v, want %v", i, decoded, uid)
		}

		// Verify timestamp preservation (to millisecond precision)
		originalTs, _ := uuidv7.TimestampFromV7(uid)
		decodedTs, _ := uuidv7.TimestampFromV7(decoded)
		
		originalTime, _ := originalTs.Time()
		decodedTime, _ := decodedTs.Time()

		if !originalTime.Equal(decodedTime) {
			t.Errorf("Timestamp mismatch at iteration %d: got %v, want %v", 
				i, decodedTime, originalTime)
		}

		// Ensure timestamps are strictly monotonic
		if i > 0 && !originalTime.After(lastTime) {
			t.Errorf("Non-monotonic timestamp at iteration %d: %v not after %v", 
				i, originalTime, lastTime)
		}
		lastTime = originalTime
	}
}

func TestRandomTemporalStressTest(t *testing.T) {
	const numOperations = 1_000_000

	// Define the valid time range for UUIDv7 (starts at Unix epoch)
	minTime := time.Unix(0, 0)                                           // Unix epoch
	maxTime := time.Date(2262, 04, 11, 23, 47, 16, 854775807, time.UTC) // Max UUIDv7 timestamp

	// Create a source of randomness
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	seenTimes := make(map[time.Time]bool)

	for i := 0; i < numOperations; i++ {
		// Generate random timestamp between min and max time
		randomUnix := rnd.Int63n(maxTime.Unix() - minTime.Unix() + 1)
		randomNano := rnd.Int63n(1_000_000_000) // Add random nanoseconds
		randomTime := time.Unix(randomUnix, randomNano)

		// Skip if we've already tested this exact timestamp
		if seenTimes[randomTime] {
			continue
		}
		seenTimes[randomTime] = true

		// Create UUID with random timestamp
		uid := uuidv7.Must(uuidv7.NewV7AtTime(randomTime))

		// Convert to short encoding and back
		encoded := ShortEncodeUUID(uid)
		decoded, err := DecodeUUID(encoded)
		if err != nil {
			t.Fatalf("Failed to decode at iteration %d (time: %v): %v", 
				i, randomTime, err)
		}

		if decoded != uid {
			t.Errorf("Round trip failed at iteration %d (time: %v): got %v, want %v", 
				i, randomTime, decoded, uid)
		}

		// Verify timestamp preservation (to millisecond precision)
		originalTs, err := uuidv7.TimestampFromV7(uid)
		if err != nil {
			t.Fatalf("Failed to extract timestamp from original UUID at iteration %d: %v", 
				i, err)
		}
		
		decodedTs, err := uuidv7.TimestampFromV7(decoded)
		if err != nil {
			t.Fatalf("Failed to extract timestamp from decoded UUID at iteration %d: %v", 
				i, err)
		}

		originalTime, err := originalTs.Time()
		if err != nil {
			t.Fatalf("Failed to convert original timestamp at iteration %d: %v", 
				i, err)
		}

		decodedTime, err := decodedTs.Time()
		if err != nil {
			t.Fatalf("Failed to convert decoded timestamp at iteration %d: %v", 
				i, err)
		}

		if !originalTime.Equal(decodedTime) {
			t.Errorf("Timestamp mismatch at iteration %d: got %v, want %v (original time: %v)", 
				i, decodedTime, originalTime, randomTime)
		}

		// Log every 100,000 operations to show progress
		if i > 0 && i%100_000 == 0 {
			t.Logf("Completed %d operations", i)
		}
	}

	t.Logf("Tested %d unique random timestamps", len(seenTimes))
}

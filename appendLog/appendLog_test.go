package appendLog

import (
	"testing"
	"time"

	uuidv7 "github.com/gofrs/uuid/v5"
)

func TestLogEntryEncode(t *testing.T) {
	fixedUUID := uuidv7.Must(uuidv7.FromString("019336d6-286a-7faf-82e2-282dc8fe40a4"))
	fixedDate := time.Date(2023, 5, 15, 12, 34, 56, 0, time.UTC)

	entry := LogEntry{
		Id:       fixedUUID,
		Type:     "INFO",
		Date:     fixedDate,
		Hostname: "example.com",
		Source:   "app.go",
		Content:  "Test message",
	}

	encoded := entry.Encode()
	expected := "019336d6-286a-7faf-82e2-282dc8fe40a4 INFO 2023-05-15T12:34:56Z example.com app.go Test message"
	if encoded != expected {
		t.Errorf("Encoded string does not match expected.\nGot: %s\nWant: %s", encoded, expected)
	}
}

func TestLogEntryDecode(t *testing.T) {
	input := "019336d6-286a-7faf-82e2-282dc8fe40a4 INFO 2023-05-15T12:34:56Z example.com app.go Test message"
	expectedUUID := uuidv7.Must(uuidv7.FromString("019336d6-286a-7faf-82e2-282dc8fe40a4"))
	expectedDate := time.Date(2023, 5, 15, 12, 34, 56, 0, time.UTC)

	entry, err := Decode(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if entry.Id != expectedUUID {
		t.Errorf("Decoded ID does not match. Got: %s, Want: %s", entry.Id, expectedUUID)
	}
	if entry.Type != "INFO" {
		t.Errorf("Decoded Type does not match. Got: %s, Want: INFO", entry.Type)
	}
	if !entry.Date.Equal(expectedDate) {
		t.Errorf("Decoded Date does not match. Got: %s, Want: %s", entry.Date, expectedDate)
	}
	if entry.Hostname != "example.com" {
		t.Errorf("Decoded Hostname does not match. Got: %s, Want: example.com", entry.Hostname)
	}
	if entry.Source != "app.go" {
		t.Errorf("Decoded Source does not match. Got: %s, Want: app.go", entry.Source)
	}
	if entry.Content != "Test message" {
		t.Errorf("Decoded Content does not match. Got: %s, Want: Test message", entry.Content)
	}
}

package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig tests the 'config' command displays configuration correctly
func TestConfig(t *testing.T) {
	cleanup, configFile, _ := setupTestEnv(t)
	defer cleanup()

	// Execute the config command
	stdout, stderr, err := runT5Command(t, "--config", configFile, "config")
	if err != nil {
		t.Fatalf("Config command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify the config file exists
	_, err = os.Stat(configFile)
	require.NoError(t, err, "Config file should exist")

	// Verify that the config output contains expected sections from the config file
	assert.Contains(t, stdout, "todo:", "Config output should contain 'todo:' section")
	assert.Contains(t, stdout, "eventstore:", "Config output should contain 'eventstore:' section")
	
	// Verify specific configuration values
	assert.Contains(t, stdout, "file: ./todo.txt", 
		"Config output should contain the correct todo file path")
	assert.Contains(t, stdout, "prefershortids: true", 
		"Config output should contain the correct ID preference")
	assert.Contains(t, stdout, "enforcecompletiondate: true", 
		"Config output should contain the correct completion date enforcement")

	t.Logf("Config command completed successfully")
}
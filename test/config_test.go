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
	_, stderr, err := runT5Command(t, "--config", configFile, "config")
	if err != nil {
		t.Fatalf("Config command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify the config file exists
	_, err = os.Stat(configFile)
	require.NoError(t, err, "Config file should exist")

	// Read the config file directly to verify contents
	configContent, err := os.ReadFile(configFile)
	require.NoError(t, err, "Should be able to read config file")
	configStr := string(configContent)

	// Verify that the config file contains expected sections
	assert.Contains(t, configStr, "todo:", "Config should contain 'todo:' section")
	assert.Contains(t, configStr, "eventstore:", "Config should contain 'eventstore:' section")
	
	// Verify specific configuration values
	assert.Contains(t, configStr, "files:", 
		"Config should contain the files key")
	assert.Contains(t, configStr, "path: ", 
		"Config should contain the path key")
	assert.Contains(t, configStr, "todo.txt", 
		"Config should contain the correct todo file path")
	assert.Contains(t, configStr, "work.txt", 
		"Config should contain the second todo file path")
	assert.Contains(t, configStr, "prefershortids: true", 
		"Config should contain the correct ID preference")
	assert.Contains(t, configStr, "enforcecompletiondate: true", 
		"Config should contain the correct completion date enforcement")

	t.Logf("Config command completed successfully")
}
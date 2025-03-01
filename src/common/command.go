package common

import (
	"fmt"
	"os/exec"
	"strings"
)

// CommandExecutor provides optimized command execution with caching
type CommandExecutor struct {
	// No fields needed for now
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

// Execute executes a command and returns its output
func (c *CommandExecutor) Execute(name string, args ...string) (string, error) {
	// Create a cache key from the command and arguments
	cacheKey := fmt.Sprintf("cmd:%s:%s", name, strings.Join(args, ":"))

	// Check if the result is in the cache
	if cachedResult, found := CommandCache.Get(cacheKey); found {
		return cachedResult.(string), nil
	}

	// Execute the command
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Cache the result
	result := strings.TrimSpace(string(output))
	CommandCache.Set(cacheKey, result)

	return result, nil
}

// ExecuteWithStdin executes a command with stdin and returns its output
func (c *CommandExecutor) ExecuteWithStdin(name string, args ...string) (string, error) {
	// Commands with stdin cannot be cached reliably
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// GlobalCommandExecutor is a global instance of CommandExecutor
var GlobalCommandExecutor = NewCommandExecutor()

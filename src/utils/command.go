package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

type CommandExecutor struct {
	// TODO: Add fields here
}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

func (c *CommandExecutor) Execute(name string, args ...string) (string, error) {
	cacheKey := fmt.Sprintf("cmd:%s:%s", name, strings.Join(args, ":"))

	if cachedResult, found := CommandCache.Get(cacheKey); found {
		return cachedResult.(string), nil
	}

	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(string(output))
	CommandCache.Set(cacheKey, result)

	return result, nil
}

func (c *CommandExecutor) ExecuteWithStdin(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

var GlobalCommandExecutor = NewCommandExecutor()

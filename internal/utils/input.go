package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ReadInput reads user input with a prompt
func ReadInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// ReadInputWithDefault reads user input with a default value
func ReadInputWithDefault(prompt, defaultValue string) string {
	input, err := ReadInput(prompt)
	if err != nil || input == "" {
		return defaultValue
	}
	return input
}

// ValidateInput validates input against a regex pattern
func ValidateInput(input, pattern string) bool {
	matched, err := regexp.MatchString(pattern, input)
	return err == nil && matched
}

// ReadValidatedInput reads and validates user input
func ReadValidatedInput(prompt, pattern, errorMsg string) (string, error) {
	for {
		input, err := ReadInput(prompt)
		if err != nil {
			return "", err
		}

		if ValidateInput(input, pattern) {
			return input, nil
		}

		PrintError(errorMsg)
	}
}

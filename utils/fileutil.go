package utils

import (
	"bufio"
	"fmt"
	"os"
)

func AppendToFile(fileName string, text string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(text + "\n")
	if err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}

	return nil
}

func ReadFromFile(fileName string) ([]string, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

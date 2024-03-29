package utils

import (
	"bufio"
	"fmt"
	"github.com/maxheckel/inchi-worm/model"
	"os"
)

func ReadFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		// handle the error here
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		// handle the error here
		return nil, err
	}
	return lines, nil
}

func WriteOutput(result []model.Inchi, name string) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	content := ""
	for _, inchi := range result {
		content += fmt.Sprintf("%s	%s\n", inchi.Key, inchi.Value)
	}
	_, err = out.WriteString(content)
	return err
}

func WriteLine(inchi model.Inchi, name string) error {

	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	line := fmt.Sprintf("%s	%s\n", inchi.Key, inchi.Value)
	if _, err = f.WriteString(line); err != nil {
		return err
	}
	return nil
}

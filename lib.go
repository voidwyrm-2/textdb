package main

import (
	"bufio"
	"os"
	"strings"
)

func ForEach[arrType, retType any](arr []arrType, fn func(v arrType, acc retType) retType) retType {
	var accum retType
	for _, item := range arr {
		accum = fn(item, accum)
	}
	return accum
}

func StartswithAny(s string, check []string) bool {
	for _, c := range check {
		if strings.HasPrefix(s, c) {
			return true
		}
	}
	return false
}

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content, nil
}

func writeFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

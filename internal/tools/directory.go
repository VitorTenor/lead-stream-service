package tools

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const projectNamespace = "lead-stream-service"

func FindProjectRoot() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if strings.HasSuffix(pwd, projectNamespace) {
			break
		}
		parentDir := filepath.Dir(pwd)
		if parentDir == pwd {
			return "", errors.New("project folder does not exist")
		}
		pwd = parentDir
	}
	return pwd, nil
}

package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func PickRandomFileFromDirectory(dirpath string) (string, error) {
	// Append "/" if path does not end with it
	if !(strings.HasSuffix(dirpath, "/") || strings.HasSuffix(dirpath, "\\")) {
		dirpath += "/"
	}

	files, err := os.ReadDir(dirpath)
	if err != nil {
		return "", fmt.Errorf("error reading %v directory contents: %s", dirpath, err)
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	filepath := dirpath + filenames[rand.Intn(len(filenames))]
	// log.Printf("[TRACE] call: PickRandomFileFromDirectory -> filepath: %v", filepath)

	return filepath, nil
}

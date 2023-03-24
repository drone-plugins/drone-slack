package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func contents(str string) (string, error) {
	// Check for the empty string
	if str == "" {
		return str, nil
	}

	isFilePath := false

	// See if the string is referencing a URL
	if u, err := url.Parse(str); err == nil {
		switch u.Scheme {
		case "http", "https":
			res, err := http.Get(str)
			if err != nil {
				return "", err
			}

			defer res.Body.Close()
			b, err := io.ReadAll(res.Body)
			if err != nil {
				return "", fmt.Errorf("could not read response: %w", err)
			}

			return string(b), nil

		case "file":
			// Fall through to file loading
			str = u.Path
			isFilePath = true
		}
	}

	// See if the string is referencing a file
	_, err := os.Stat(str)
	if err == nil {
		b, err := os.ReadFile(str)
		if err != nil {
			return "", fmt.Errorf("could not load file %s: %w", str, err)
		}

		return string(b), nil
	}

	if isFilePath {
		return "", fmt.Errorf("could not load file %s: %w", str, err)
	}

	// Its a regular string
	return str, nil
}

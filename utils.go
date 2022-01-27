package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	AMD64_IDENTIFIER = "amd64"
	X64_IDENTIFER    = "x64"
	X86_64_IDENTIFER = "x86_64"
)

// setup global logging configuration and make logging level
// configurable using LOG_LEVEL environment variable
func setupLogging() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "debug"
	}
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	log.SetOutput(os.Stdout)
	log.SetLevel(ll)
}

// download file located at the given URL to filepath
func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// generate universal binary asset name for github based on arch
func generateUniversalAssetName(amd64AssetName string, identifier string) string {
	var architecture string

	// if the asset name contains amd64, x64 or x86_64, convert
	// that architecture identifier to the universal identifier.
	// assumes that the architecture is x86_64 if nothing is provided.
	if strings.Contains(amd64AssetName, AMD64_IDENTIFIER) {
		architecture = AMD64_IDENTIFIER
	} else if strings.Contains(amd64AssetName, X64_IDENTIFER) {
		architecture = X64_IDENTIFER
	} else {
		architecture = X86_64_IDENTIFER
	}

	return strings.ReplaceAll(amd64AssetName, architecture, identifier)
}

func findBinaryPath(directory string, name string) (string, error) {
	var binaryPath string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) == name {
			binaryPath = path
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return binaryPath, nil
}

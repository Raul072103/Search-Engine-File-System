package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// CalculateFileMD5 Calculates the MD5 hash of a file
func CalculateFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes), nil
}

// CalculateDirectoryMD5 Calculates the MD5 hash of a directory (only the first-level items)
func CalculateDirectoryMD5(path string) (string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer dir.Close()

	entries, err := dir.Readdir(0)
	if err != nil {
		return "", err
	}

	// Important, we want consistency
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	hash := md5.New()

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			_, err := io.WriteString(hash, entryPath)
			if err != nil {
				return "", err
			}
		} else {
			fileHash, err := CalculateFileMD5(entryPath)
			if err != nil {
				return "", err
			}

			_, err = io.WriteString(hash, entryPath+fileHash)
			if err != nil {
				return "", err
			}
		}
	}

	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes), nil
}

package sha

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func SHA256FromFile(path string) (string, error) {
	openFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer openFile.Close()
	s := sha256.New()
	_, err = io.Copy(s, openFile)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(s.Sum(nil)), nil
}

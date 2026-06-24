package sha

import (
	"log"
	"testing"
)

const trueHash = "1a5376ad727d65213a79f3108541cf95012969a0d3064f108b5dd6e7f8c19b89"

func TestSHA256FromFile(t *testing.T) {
	hash, err := SHA256FromFile("test_data/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	if hash != trueHash {
		log.Fatalf("expected %s, got %s", trueHash, hash)
	}
	log.Printf("SHA256 of test_data/test.txt: %s", hash)
}

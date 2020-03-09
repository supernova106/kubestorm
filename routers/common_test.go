package routers

import "testing"

// TODO: unittest
func TestWriteToFile(t *testing.T) {
	err := writeToFile([]byte("hello"), "/tmp/kubestorm-test-file")
	if err != nil {
		t.Errorf("writeToFile([]byte(\"hello\"), \"/tmp/kubestorm-test-file\") failed, %v", err)
	}
}

func TestHash(t *testing.T) {
	newHash := hash("hello")
	if newHash != "5d41402abc4b2a76b9719d911017c592" {
		t.Errorf("hash(\"hello\") failed, expected %v", newHash)
	}
}

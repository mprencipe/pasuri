package main

import "testing"

func TestMakeHash(t *testing.T) {
	correctHash := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	hash := MakeHash("test")
	if hash != correctHash {
		t.Errorf("Hash was incorrect, got: %s, wanted: %s.", hash, correctHash)
	}
}

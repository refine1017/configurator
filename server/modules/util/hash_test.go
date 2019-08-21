package util

import (
	"testing"
)

func Test_HashGet(t *testing.T) {
	h1 := NewHasher(HASH_LETTER)
	if len(h1.Get(0)) != 0 {
		t.Errorf("HASH_LETTER Get 0 != 0")
	}
	if len(h1.Get(5)) != 5 {
		t.Errorf("HASH_LETTER Get 5 != 5")
	}

	h2 := NewHasher(HASH_LETTER_AND_NUMBER)
	if len(h2.Get(0)) != 0 {
		t.Errorf("HASH_LETTER_AND_NUMBER Get 0 != 0")
	}
	if len(h2.Get(5)) != 5 {
		t.Errorf("HASH_LETTER_AND_NUMBER Get 5 != 5")
	}
}

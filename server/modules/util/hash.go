package util

import (
	"server/modules/algorithm"
	"strings"
)

const (
	HASH_LETTER            = iota // 单字母
	HASH_NUMBER                   // 数字
	HASH_LETTER_AND_NUMBER        // 字母加数字
)

var defaultHasher *Hasher

// 目录Hash
type Hasher struct {
	data []string
}

func NewHasher(catalog int) *Hasher {
	hasher := &Hasher{}

	switch catalog {
	case HASH_LETTER:
		hasher.letterData()
	case HASH_NUMBER:
		hasher.numberData()
	case HASH_LETTER_AND_NUMBER:
		hasher.letterAndNumberData()
	}

	return hasher
}

func (h *Hasher) letterData() {
	h.data = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z",
	}
}

func (h *Hasher) numberData() {
	h.data = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
}

func (h *Hasher) letterAndNumberData() {
	h.data = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
}

func (h *Hasher) Get(n int) string {
	if len(h.data) == 0 || n <= 0 {
		return ""
	}

	var strs = make([]string, n)

	for i := 0; i < n; i++ {
		index := algorithm.Random(0, int64(len(h.data)-1))

		strs[i] = h.data[index]
	}

	return strings.Join(strs, "")
}

func HashGet(n int) string {
	return defaultHasher.Get(n)
}

func init() {
	defaultHasher = NewHasher(HASH_LETTER_AND_NUMBER)
}

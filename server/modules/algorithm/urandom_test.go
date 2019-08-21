package algorithm

import (
	"testing"
)

func Test_Random(t *testing.T) {
	if Random(0, 0) != 0 {
		t.Errorf("Random [0,0] != 0")
	}

	if Random(1, 0) != 0 {
		t.Errorf("Random [1,0] != 0")
	}

	if Random(0, 1) > 1 {
		t.Errorf("Random [0,1] != 1")
	}
}

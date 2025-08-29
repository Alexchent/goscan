package help

import "testing"

func TestScale(t *testing.T) {
	err := Scale("/Users/chentao02/Downloads")
	if err != nil {
		t.Error(err)
	}
}

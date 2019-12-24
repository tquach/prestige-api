package common

import "testing"

func TestString(t *testing.T) {
	if Pending.String() != "Pending" {
		t.Logf("Expected %s but got %s", "Pending", Pending)
		t.Fail()
	}
}

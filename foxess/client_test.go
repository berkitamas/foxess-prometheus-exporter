package foxess

import (
	"testing"
)

func TestCalculateSignature(t *testing.T) {
	// Test vector:
	//   calculateSignature("/op/v0/device/list", "abcdefghij012345689", 1705809089)
	//   → "68a007c2450d6697fbe2990f92000269"
	path := "/op/v0/device/list"
	apiKey := "abcdefghij012345689"
	timestamp := int64(1705809089)
	expected := "68a007c2450d6697fbe2990f92000269"

	got := calculateSignature(path, apiKey, timestamp)
	if got != expected {
		t.Errorf("calculateSignature() = %v, want %v", got, expected)
	}
}

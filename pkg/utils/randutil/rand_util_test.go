package randutil

import (
	"slices"
	"testing"
)

// TestRandString 测试 RandString
func TestRandString(t *testing.T) {
	history := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		ret := RandString(AlphabetLowerAlphaNumeric, 16)
		if slices.Contains(history, ret) {
			t.Fatalf("repeated results: %q", ret)
		}
		t.Logf("result: %q", ret)
		history = append(history, ret)
	}
}

// BenchmarkRandString_36x16 RandString 基准测试， alphabet 长度为 32 ，生成长度 16
func BenchmarkRandString_36x16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandString(AlphabetLowerAlphaNumeric, 16)
	}
}

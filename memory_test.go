package main

import (
	"testing"
)

func TestSampleMemorySaturation(t *testing.T) {
	ms, err := SampleMemorySaturation()
	if err != nil {
		t.Error(err)
	}

	if ms.SwapUsage == float64(0) {
		t.Error("ms should not be 0, got 0")
	}
}

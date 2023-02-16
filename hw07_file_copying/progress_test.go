package main

import (
	"testing"
	"time"
)

func Test_Run(t *testing.T) {
	bar := CreateNew(5)
	bar.Start()
	bar.Prefix("Increment values:")
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 50)
		bar.Increment()
	}
	if actual := bar.Get(); actual != 5 {
		t.Errorf("Expected: %d; actual: %d", 5, actual)
	}
	bar.Finish()
}

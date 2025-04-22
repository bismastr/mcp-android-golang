package adb_test

import (
	"testing"

	"github.com/bismastr/mcp-android-automation/internal/adb"
)

func Test(t *testing.T) {
	device, err := adb.NewAdbDevice()
	if err != nil {
		t.Fatalf("Error creating device: %v", err)
	}

	device.GetDevice()
	device.TakeScreenshot()
}

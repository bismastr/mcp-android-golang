package adb_test

import (
	"testing"

	"github.com/bismastr/mcp-android-automation/internal/adb"
)

func TestScreenshot(t *testing.T) {
	device, err := adb.NewAdbDevice()
	if err != nil {
		t.Fatalf("Error creating device: %v", err)
	}

	device.TakeScreenshotBase64()
}

func TestGetHeirarchy(t *testing.T) {
	device, err := adb.NewAdbDevice()
	if err != nil {
		t.Fatalf("Error creating device: %v", err)
	}

	result, err := device.GetUIHierarchy()
	if err != nil {
		t.Fatalf("Error get ui: %v", err)
	}

	t.Log(result)
}

func TestXMLParser(t *testing.T) {
	device, err := adb.NewAdbDevice()
	if err != nil {
		t.Fatalf("Error creating device: %v", err)
	}

	result, err := device.GetUIHierarchy()
	if err != nil {
		t.Fatalf("Error get ui: %v", err)
	}

	t.Log(result)
}

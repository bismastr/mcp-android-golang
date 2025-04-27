package adb

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"time"

	"github.com/electricbubble/gadb"
	"golang.org/x/image/draw"
)

type AndroidDevice struct {
	adb gadb.Device
}

func NewAdbDevice() (*AndroidDevice, error) {
	adbClient, err := gadb.NewClient()
	if err != nil {
		return nil, err
	}

	device, err := adbClient.DeviceList()
	if err != nil {
		return &AndroidDevice{
			adb: gadb.Device{},
		}, err
	}

	return &AndroidDevice{
		adb: device[0],
	}, nil
}

func (a *AndroidDevice) GetDevice() {
	fmt.Println(a.adb.DeviceInfo())
}

func (d *AndroidDevice) ShellCommand(cmd string, args ...string) (string, error) {
	return d.adb.RunShellCommand(cmd, args...)
}

func (d *AndroidDevice) Tap(x, y int) error {
	cmd := fmt.Sprintf("input tap %d %d", x, y)
	_, err := d.ShellCommand(cmd)
	return err
}

func (d *AndroidDevice) TakeScreenshot() (*os.File, error) {
	// Use absolute path within workspace
	dir := "/Users/bytedance/Documents/repo/mcp-android/screenshots"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")
	screenshotPath := fmt.Sprintf("%s/%s.jpg", dir, timestamp) // Changed to jpg

	devicePath := "/sdcard/screenshot.png"
	_, err = d.ShellCommand("screencap", "-p", devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %v", err)
	}

	// Create a temporary file to store the pulled PNG
	tempFile, err := os.CreateTemp(dir, "temp_*.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath) // Clean up temp file when done

	// Pull the screenshot from device
	err = d.adb.Pull(devicePath, tempFile)
	if err != nil {
		tempFile.Close()
		return nil, fmt.Errorf("failed to pull screenshot: %v", err)
	}
	tempFile.Close()

	// Clean up the screenshot on device
	_, err = d.ShellCommand("rm", devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to remove device screenshot: %v", err)
	}

	// Open the temp file for reading
	tempFile, err = os.Open(tempFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open temp file: %w", err)
	}
	defer tempFile.Close()

	// Decode PNG
	img, err := png.Decode(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a new image with half the dimensions
	newWidth := width / 2
	newHeight := height / 2
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Scale down the image
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Create the output JPEG file
	outFile, err := os.Create(screenshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	// Encode as JPEG with 75% quality
	err = jpeg.Encode(outFile, dst, &jpeg.Options{Quality: 75})
	if err != nil {
		outFile.Close()
		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
	}

	// Reset file pointer to beginning
	if _, err := outFile.Seek(0, 0); err != nil {
		outFile.Close()
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	return outFile, nil
}

// TakeScreenshotBase64 takes a screenshot and returns it as a base64 encoded string
func (d *AndroidDevice) TakeScreenshotBase64() (string, error) {
	// Use a buffer instead of a file
	var buf bytes.Buffer

	devicePath := "/sdcard/screenshot.png"
	_, err := d.ShellCommand("screencap", "-p", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to capture screenshot: %v", err)
	}

	// Create a temporary file to store the pulled PNG
	dir := "/Users/bytedance/Documents/repo/mcp-android/screenshots"
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	tempFile, err := os.CreateTemp(dir, "temp_*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath) // Clean up temp file when done

	// Pull the screenshot from device
	err = d.adb.Pull(devicePath, tempFile)
	if err != nil {
		tempFile.Close()
		return "", fmt.Errorf("failed to pull screenshot: %v", err)
	}
	tempFile.Close()

	// Clean up the screenshot on device
	_, err = d.ShellCommand("rm", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to remove device screenshot: %v", err)
	}

	// Open the temp file for reading
	tempFile, err = os.Open(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open temp file: %w", err)
	}
	defer tempFile.Close()

	// Decode PNG
	img, err := png.Decode(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a new image with half the dimensions
	newWidth := width / 2
	newHeight := height / 2
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Scale down the image
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Encode as JPEG with 75% quality directly to buffer
	err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 75})
	if err != nil {
		return "", fmt.Errorf("failed to encode JPEG: %w", err)
	}

	// Convert to base64
	base64Data := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(base64Data, buf.Bytes())

	return string(base64Data), nil
}

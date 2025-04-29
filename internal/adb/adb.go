package adb

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/beevik/etree"
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

func (d *AndroidDevice) GetUIHierarchy() (string, error) {
	devicePath := "/sdcard/window_dump.xml"

	_, err := d.ShellCommand("uiautomator", "dump", "--compressed", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to dump UI hierarchy: %v", err)
	}

	xmlContent, err := d.ShellCommand("cat", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read UI hierarchy file: %v", err)
	}

	_, err = d.ShellCommand("rm", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to remove temporary file: %v", err)
	}

	return xmlContent, nil
}

func (d *AndroidDevice) ParseXML(xmlData string) ([]string, error) {
	doc := etree.NewDocument()

	err := doc.ReadFromString(xmlData)
	if err != nil {
		return nil, err
	}

	root := doc.SelectElement("hierarchy")
	if root == nil {
		return nil, fmt.Errorf("hierarchy element not found in XML")
	}

	result := []string{}
	for _, e := range root.FindElements("//node") {
		bounds := e.SelectAttr("bounds")
		if bounds != nil {
			result = append(result, bounds.Value)
		}
	}

	return result, nil
}

func (d *AndroidDevice) TakeScreenshotBase64() (string, error) {
	var buf bytes.Buffer

	devicePath := "/sdcard/screenshot.png"
	_, err := d.ShellCommand("screencap", "-p", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to capture screenshot: %v", err)
	}

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
	defer os.Remove(tempFilePath)

	err = d.adb.Pull(devicePath, tempFile)
	if err != nil {
		tempFile.Close()
		return "", fmt.Errorf("failed to pull screenshot: %v", err)
	}
	tempFile.Close()

	_, err = d.ShellCommand("rm", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to remove device screenshot: %v", err)
	}

	tempFile, err = os.Open(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open temp file: %w", err)
	}
	defer tempFile.Close()

	img, err := png.Decode(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newWidth := width / 2
	newHeight := height / 2
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 75})
	if err != nil {
		return "", fmt.Errorf("failed to encode JPEG: %w", err)
	}

	base64Data := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(base64Data, buf.Bytes())

	return string(base64Data), nil
}

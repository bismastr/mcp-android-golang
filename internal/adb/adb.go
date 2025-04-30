package adb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"regexp"
	"strconv"

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

func (d *AndroidDevice) SendKeys(text string) error {
	cmd := fmt.Sprintf("input text \"%s\"", text)
	_, err := d.ShellCommand(cmd)
	if err != nil {
		return err
	}

	return err
}

func (d *AndroidDevice) PressEnter() error {
	cmd := "input keyevent 66"
	_, err := d.ShellCommand(cmd)
	return err
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

	parsed, err := d.ParseXML(xmlContent)
	if err != nil {
		return "", err
	}

	elements, err := d.CollectElements(parsed.Nodes)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(elements)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (d *AndroidDevice) ParseXML(xmlData string) (*Hierarchy, error) {

	var hierarchy Hierarchy
	err := xml.Unmarshal([]byte(xmlData), &hierarchy)
	if err != nil {
		fmt.Printf("Error parsing XML: %v\n", err)
		return nil, err
	}

	return &hierarchy, nil
}

func (d *AndroidDevice) CollectElements(nodes []Node) ([]UIElement, error) {
	var result []UIElement
	for _, node := range nodes {
		x, y, w, h, err := parseBounds(node.Bounds)
		if err != nil {
			fmt.Printf("Skipping invalid node: %v\n", err)
			continue
		}

		temp := UIElement{
			Text:        node.Text,
			Class:       node.Class,
			ContentDesc: node.ContentDesc,
			ResourceID:  node.ResourceID,
			Bounds:      node.Bounds,
			X:           x,
			Y:           y,
			Width:       w,
			Height:      h,
		}

		result = append(result, temp)

		childElements, err := d.CollectElements(node.ChildNodes)
		if err != nil {
			return nil, err
		}

		result = append(result, childElements...)
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

func parseBounds(bounds string) (x, y, width, height int, err error) {
	re := regexp.MustCompile(`\[(\d+),(\d+)\]\[(\d+),(\d+)\]`)
	matches := re.FindStringSubmatch(bounds)

	if len(matches) != 5 {
		return 0, 0, 0, 0, fmt.Errorf("invalid bounds format: %s", bounds)
	}

	left, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid left coordinate: %w", err)
	}

	top, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid top coordinate: %w", err)
	}

	right, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid right coordinate: %w", err)
	}

	bottom, err := strconv.Atoi(matches[4])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid bottom coordinate: %w", err)
	}

	if right < left || bottom < top {
		return 0, 0, 0, 0, fmt.Errorf("invalid rectangle coordinates: %s", bounds)
	}

	x = (left + right) / 2
	y = (top + bottom) / 2
	width = right - left
	height = bottom - top
	return
}

package adb

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/electricbubble/gadb"
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
	wd, _ := os.Getwd()
	dir := path.Join(wd, "screenshot")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")
	screenshotPath := fmt.Sprintf("%s/%s.png", dir, timestamp)

	devicePath := "/sdcard/screenshot.png"
	_, err = d.ShellCommand("screencap", "-p", devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %v", err)
	}

	file, err := os.Create(screenshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local file: %w", err)
	}

	err = d.adb.Pull(devicePath, file)
	if err != nil {
		return nil, fmt.Errorf("failed to pull screenshot: %v", err)
	}

	_, err = d.ShellCommand("rm", devicePath)

	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	return file, nil
}

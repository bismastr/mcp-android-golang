package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/bismastr/mcp-android-automation/internal/adb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	device, err := adb.NewAdbDevice()
	if err != nil {
		log.Fatalf("Failed to initialize ADB device: %v", err)
	}

	s := server.NewMCPServer(
		"mcp-android-adb",
		"1.0.0",
		server.WithLogging(),
	)

	s.AddTool(mcp.NewTool("Take a screenshot and get base64 data"),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			file, err := device.TakeScreenshot()
			if err != nil {
				return nil, fmt.Errorf("failed to take screenshot: %w", err)
			}

			defer file.Close()

			// Read the file content
			fileInfo, err := file.Stat()
			if err != nil {
				return nil, fmt.Errorf("failed to get file info: %w", err)
			}

			// Read the file content
			fileInfo, err = file.Stat()
			if err != nil {
				return nil, fmt.Errorf("failed to get file info: %w", err)
			}

			fileContent := make([]byte, fileInfo.Size())
			_, err = file.Read(fileContent)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}

			// Convert to base64 for transmission
			base64Data := base64.StdEncoding.EncodeToString(fileContent)

			// Create image content using the helper function
			imageContent := mcp.NewToolResultImage("Screenshot", base64Data, "image/png")

			// Return the result with file content and metadata
			return imageContent, nil
		},
	)

	fmt.Println("Server is running")
	err = server.ServeStdio(s)
	if err != nil {
		log.Fatalf("Error serving stdio %v", err)
	}

}

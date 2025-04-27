package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	"github.com/bismastr/mcp-android-automation/internal/adb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	sseMode := flag.Bool("sse", false, "Run in SSE mode instead of stdio mode")
	flag.Parse()

	device, err := adb.NewAdbDevice()
	if err != nil {
		log.Fatalf("Failed to initialize ADB device: %v", err)
	}

	s := server.NewMCPServer(
		"mcp-android-adb",
		"1.0.0",
		server.WithLogging(),
	)

	s.AddTool(mcp.NewTool("take-sceenshot",
		mcp.WithDescription(
			"Take a screenshot of the mobile device. Use this to understand what's on screen, if you need to press an element that is available through view hierarchy then you must list elements on screen instead. Do not cache this result.",
		),
	),
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

	// Run server in appropriate mode
	if *sseMode {
		sseServer := server.NewSSEServer(s)
		log.Printf("Starting SSE server on localhost:8080")
		if err := sseServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}

}

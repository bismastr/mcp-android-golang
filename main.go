package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/bismastr/mcp-android-automation/internal/adb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
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
			base64Data, err := device.TakeScreenshotBase64()
			if err != nil {
				return nil, fmt.Errorf("failed to take screenshot: %w", err)
			}

			imageContent := mcp.NewToolResultImage("Screenshot", base64Data, "image/jpeg")

			return imageContent, nil
		},
	)

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

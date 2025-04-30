package main

import (
	"flag"
	"log"

	"github.com/bismastr/mcp-android-automation/internal/adb"
	"github.com/bismastr/mcp-android-automation/internal/tools"
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

	addTools(s, device)

	if *sseMode {
		sseServer := server.NewSSEServer(s)
		log.Printf("Starting SSE server on localhost")
		if err := sseServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

func addTools(s *server.MCPServer, d *adb.AndroidDevice) {
	tools := []func(*server.MCPServer, *adb.AndroidDevice){
		tools.AddToolListelement,
		tools.AddToolScreenshot,
		tools.AddToolTapWithCoordinate,
		tools.AddToolSendKeys,
	}

	for _, tool := range tools {
		tool(s, d)
	}
}

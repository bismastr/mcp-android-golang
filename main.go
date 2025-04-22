package main

import "github.com/mark3labs/mcp-go/server"

func main() {
	server.NewMCPServer(
		"mcp-android-adb",
		"1.0.0",
		server.WithLogging(),
	)

}

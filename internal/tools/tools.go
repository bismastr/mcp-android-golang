package tools

import (
	"context"
	"fmt"

	"github.com/bismastr/mcp-android-automation/internal/adb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddToolScreenshot(s *server.MCPServer, device *adb.AndroidDevice) {
	s.AddTool(mcp.NewTool("take-sceenshot",
		mcp.WithDescription(
			"Take a screenshot of the mobile device. Use this to understand what's on screen, if you need to press an element that is available then you must list elements on screen instead. Do not cache this result.",
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
}

func AddToolListelement(s *server.MCPServer, device *adb.AndroidDevice) {
	s.AddTool(mcp.NewTool("list-element", mcp.WithDescription("List elements on screen and their coordinates, with display text, id, class and bounds. Do not cache this result.")),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			elements, err := device.GetUIHierarchy()
			if err != nil {
				return nil, err
			}

			result := mcp.NewToolResultText(elements)
			return result, nil
		},
	)
}

func AddToolTapWithCoordinate(s *server.MCPServer, device *adb.AndroidDevice) {
	s.AddTool(mcp.NewTool("tap-with-coordinate",
		mcp.WithDescription("Tap on the screen at given x,y coordinates"),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("X coordinate of the tap position"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Y coordinate of the tap position"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		x := int(request.Params.Arguments["x"].(float64))
		y := int(request.Params.Arguments["y"].(float64))

		err := device.Tap(x, y)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(fmt.Sprintf("Tap operation performed at coordinates (%d, %d)", x, y)), nil
	})
}

func AddToolSendKeys(s *server.MCPServer, device *adb.AndroidDevice) {
	s.AddTool(mcp.NewTool("input-text",
		mcp.WithDescription("Input text on Android device"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text that will be inputted"),
		),
		mcp.WithBoolean("submit",
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		err := device.SendKeys(request.Params.Arguments["text"].(string))
		if err != nil {
			return nil, err
		}

		if request.Params.Arguments["submit"].(bool) {
			device.PressEnter()
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully input text: %s", request.Params.Arguments["text"].(string))), nil
	})
}

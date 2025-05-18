# Android MCP Golang
An MCP Server that integrates the ADB command, providing device actions capabilities into LLM 

## Tools
- **input-text**
	- will execute send keys and click submit if necessary

- **tap-with-coordinate**
	- Tap on the screen at given x,y coordinates

- **list-element**
	- List elements on screen and their coordinates, with display text, id, class and bounds. Do not cache this result.

- **take-sceenshot**
	- Take a screenshot of the mobile device.

## Demo
![Screen Recording 2025-05-04 at 19 50 11](https://github.com/user-attachments/assets/f4898430-ff78-4a68-abe7-8d99c05156f3)

## Run Local

#### Clone the project

```bash
  git clone https://github.com/bismastr/mcp-android-golang.git
```

#### Go to the project directory

```bash
cd mcp-android-golang
```

#### Install dependencies

```go
go install
```

#### Open up claude and Go to Settings -> Developer -> Edit Config 
![Screenshot 2025-05-18 at 13 47 31](https://github.com/user-attachments/assets/768c2316-c4ff-48ab-a5a6-e8164e36993c)


#### Open claude_desktop_config.json and Insert your installed mcp server.

```json
{
  "mcpServers": {
    "mobile-mcp": {
      "command": "/Users/bismastr/go/bin/mcp-android-automation"
    }
  }
}
```
#### Open Up Android Simulator 
Recommended using Android Studio Built in Simulator

#### Restart Claude Desktop by Quitting the App
Try to restart the client




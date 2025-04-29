package adb

import "encoding/xml"

type Hierarchy struct {
	XMLName  xml.Name `xml:"hierarchy"`
	Rotation string   `xml:"rotation,attr"`
	Nodes    []Node   `xml:"node"`
}

type Node struct {
	Index       string `xml:"index,attr"`
	Text        string `xml:"text,attr"`
	ResourceID  string `xml:"resource-id,attr"`
	Class       string `xml:"class,attr"`
	Package     string `xml:"package,attr"`
	ContentDesc string `xml:"content-desc,attr"`
	Clickable   string `xml:"clickable,attr"`
	Bounds      string `xml:"bounds,attr"`
	ChildNodes  []Node `xml:"node"`
}

type UIElement struct {
	Text        string `json:"text"`
	Class       string `json:"class"`
	ResourceID  string `json:"resourceId"`
	ContentDesc string `json:"contentDesc"`
	Bounds      string `json:"bounds"`
	Clickable   bool   `json:"clickable"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

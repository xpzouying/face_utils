package main

// FaceInfo: information for face
// Image: head rect image in base64 encoding
// Rect: face rectangle
type FaceInfo struct {
	Image string `json:"image"`
	Rect  [4]int `json:"rect"`
}

type DetectResult struct {
	TimeUsed  float64    `json:"time_used"`
	FacesInfo []FaceInfo `json:"faces"`
}

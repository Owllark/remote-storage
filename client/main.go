package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"remote-storage/schemas"
	"time"
)

const (
	total   = 100
	barSize = 20
)

func main() {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	data, _ := json.Marshal(schemas.CdRequest{Path: "test_dir"})
	req, _ := http.NewRequest("PUT", "/cd", bytes.NewReader(data))
	client.Do(req)
	for i := 0; i <= total; i++ {
		progress := i * barSize / total
		fmt.Printf("\r[%s%s] %d%%", getProgressBar(progress), getEmptyBar(barSize-progress), i)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nTask completed!")
}

func getProgressBar(progress int) string {
	bar := ""
	for i := 0; i < progress; i++ {
		bar += "="
	}
	return bar
}

func getEmptyBar(emptyCount int) string {
	empty := ""
	for i := 0; i < emptyCount; i++ {
		empty += " "
	}
	return empty
}

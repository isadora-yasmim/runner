package cmd

import (
	"fmt"
	"net/http"
	"time"
)

func isServerRunning(port int) bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:%d/health", port)

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
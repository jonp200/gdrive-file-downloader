package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// downloadFile downloads a file from the given URL and saves it to the specified file path.
func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// getDirectDownloadLink converts a Google Drive sharing link to a direct download link.
func getDirectDownloadLink(driveLink string) (string, error) {
	// Check if the link is a valid Google Drive link
	if !strings.Contains(driveLink, "drive.google.com") {
		return "", fmt.Errorf("invalid Google Drive link")
	}

	// Extract the file ID from the link
	var fileID string
	if strings.Contains(driveLink, "/file/d/") {
		parts := strings.Split(driveLink, "/file/d/")
		if len(parts) > 1 {
			fileID = strings.Split(parts[1], "/")[0]
		}
	} else if strings.Contains(driveLink, "id=") {
		parts := strings.Split(driveLink, "id=")
		if len(parts) > 1 {
			fileID = strings.Split(parts[1], "&")[0]
		}
	} else {
		return "", fmt.Errorf("could not extract file ID from link")
	}

	if fileID == "" {
		return "", fmt.Errorf("file ID is empty")
	}

	// Create the direct download link
	directLink := fmt.
		Sprintf(
			"https://drive.usercontent.google.com/download?id=%s&export=download&authuser=0&confirm=t",
			fileID,
		)
	return directLink, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <Google Drive link> <output file path>")
		return
	}

	driveLink := os.Args[1]
	outputPath := os.Args[2]

	directLink, err := getDirectDownloadLink(driveLink)
	if err != nil {
		log.Fatalf("Failed to get direct download link: %v", err)
	}

	fmt.Printf("Downloading file from: %s\n", directLink)
	if err := downloadFile(directLink, outputPath); err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}

	fmt.Printf("File downloaded successfully to: %s\n", outputPath)
}

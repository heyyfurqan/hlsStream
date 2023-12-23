package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/heyyfurqan/hlsStream/audio"
	"github.com/heyyfurqan/hlsStream/video"
)

func main() {
	// Define your input file, output directory, resolutions, and bitrates
	resolutions := []string{"240p", "360p", "480p", "720p"}
	bitrates := []string{"400000", "800000", "1500000", "3500000"}

	audioFiles, err := os.ReadDir("audio/storage")
	if err != nil {
		log.Fatalf("Error reading audio files: %v", err)
	}

	videoFiles, err := os.ReadDir("video/storage")
	if err != nil {
		log.Fatalf("Error reading video files: %v", err)
	}

	// Segment and create playlist files for each audio file
	if len(audioFiles) == 0 {
		log.Println("No audio files found.")
	} else {
		for _, audioFile := range audioFiles {
			inputFile := filepath.Join("audio/storage", audioFile.Name())
			outputDirectory := "audio/segments"

			// Create audio segments
			audio.CreateSegments(inputFile, outputDirectory, 6)
		}
	}

	// Segment and create playlist files for each video file
	if len(videoFiles) == 0 {
		log.Println("No video files found.")
	} else {
		for _, videoFile := range videoFiles {
			inputFile := filepath.Join("video/storage", videoFile.Name())
			outputDirectory := "video/segments"

			// Create video segments
			video.CreateVideoSegments(inputFile, outputDirectory, resolutions, bitrates)

			// Create master playlist
			video.CreateMasterPlaylist(outputDirectory, resolutions, bitrates)
		}
	}

	// Start the server
	StartServer()
}

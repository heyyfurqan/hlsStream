package audio

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func CreateSegments(inputFile string, outputDirectory string, segmentDuration int) {
	// Get the base name of the input file
	baseName := filepath.Base(inputFile)
	// Remove the extension
	baseName = baseName[:len(baseName)-len(filepath.Ext(baseName))]

	// Create a new directory for this file's segments and playlist
	fileOutputDirectory := filepath.Join(outputDirectory, baseName)
	err := os.MkdirAll(fileOutputDirectory, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory %s: %v", fileOutputDirectory, err)
	}

	// Construct the FFmpeg command to segment the media file
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:a", "aac",
		"-b:a", "128k",
		"-map", "0",
		"-f", "segment",
		"-segment_time", strconv.Itoa(segmentDuration),
		"-segment_format", "mpegts",
		"-segment_list", fileOutputDirectory+"/playlist.m3u8",
		fileOutputDirectory+"/segment%03d.ts",
	)

	// Run the FFmpeg command
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

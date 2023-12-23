package video

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateVideoSegments(inputFile string, outputDirectory string, resolutions []string, bitrates []string) {
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

	for i, resolution := range resolutions {
		bitrate := bitrates[i]

		// Create a new directory for this resolution
		resolutionOutputDirectory := filepath.Join(fileOutputDirectory, "stream_"+resolution)
		err := os.MkdirAll(resolutionOutputDirectory, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", resolutionOutputDirectory, err)
		}

		// Construct the FFmpeg command to segment the video file
		cmd := exec.Command(
			"ffmpeg",
			"-i", inputFile,
			"-vf", "scale=-1:"+resolution,
			"-c:v", "libx264",
			"-b:v", bitrate,
			"-map", "0",
			"-f", "segment",
			"-segment_time", "6",
			"-segment_format", "mpegts",
			"-segment_list", resolutionOutputDirectory+"/playlist.m3u8",
			resolutionOutputDirectory+"/segment_%03d.ts",
		)

		// Run the FFmpeg command
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

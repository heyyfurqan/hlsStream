package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/grafov/m3u8"
)

func main() {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %s", err)
	}

	hlsDirectory := currentDirectory

	fileServer := http.FileServer(http.Dir(hlsDirectory))

	segmentVideos()

	createPlaylists()

	http.Handle("/", corsMiddleware(fileServer))

	port := 8080

	fmt.Printf("Serving on HTTP port: %d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func segmentVideos() {
	videoDir := "./videos"
	files, err := os.ReadDir(videoDir)
	if err != nil {
		log.Fatalf("Error reading video directory: %s", err)
	}

	for i := 1; i <= len(files); i++ {
		file := files[i-1]
		if !file.IsDir() {
			videoPath := filepath.Join(videoDir, file.Name())
			outputPath := filepath.Join("./segmented_videos", fmt.Sprintf("stream_%d", i))
			os.MkdirAll(outputPath, os.ModePerm)

			cmd := exec.Command("ffmpeg", "-i", videoPath, "-filter_complex", "[0:v]split=4[v1][v2][v3][v4]; [v1]scale=w=1280:h=720[v1out]; [v2]scale=w=854:h=480[v2out]; [v3]scale=w=640:h=360[v3out]; [v4]scale=w=426:h=240[v4out]",
				"-map", "[v1out]", "-c:v:0", "libx264", "-b:v:0", "3500k", "-map", "[v2out]", "-c:v:1", "libx264", "-b:v:1", "1500k", "-map", "[v3out]", "-c:v:2", "libx264", "-b:v:2", "800k", "-map", "[v4out]", "-c:v:3", "libx264", "-b:v:3", "400k",
				"-f", "hls", "-hls_time", "6", "-hls_playlist_type", "vod", "-hls_flags", "independent_segments", "-hls_segment_type", "mpegts", "-hls_segment_filename", outputPath+"/data%02d.ts", "-var_stream_map", "v:0 v:1 v:2 v:3", "-master_pl_name", "master.m3u8", outputPath+"/stream.m3u8")

			fmt.Printf("Running command: %v\n", cmd)

			err := cmd.Run()
			if err != nil {
				log.Fatalf("FFmpeg error for file %s: %s", file.Name(), err)
			}
		}
	}
}

func createPlaylists() {
	hlsDir := "./segmented_videos"
	videoDirs, err := os.ReadDir(hlsDir)
	if err != nil {
		log.Fatalf("Error reading HLS directory: %s", err)
	}

	for _, videoDir := range videoDirs {
		if videoDir.IsDir() {
			videoPath := filepath.Join(hlsDir, videoDir.Name())
			master := m3u8.NewMasterPlaylist()

			resolutions := []struct {
				width, height int
				bandwidth     int
			}{
				{1280, 720, 3500},
				{854, 480, 1500},
				{640, 360, 800},
				{426, 240, 400},
			}

			for i, resolution := range resolutions {
				variantParams := m3u8.VariantParams{
					Bandwidth: uint32(resolution.bandwidth),
				}

				uri := fmt.Sprintf("stream%d/stream.m3u8", i+1)
				mediaPlaylist, err := m3u8.NewMediaPlaylist(100, 100)
				if err != nil {
					log.Fatalf("Error creating media playlist: %s", err)
				}

				master.Append(uri, mediaPlaylist, variantParams)
			}

			masterBuf := master.Encode()

			masterPath := filepath.Join(videoPath, "master.m3u8")
			masterFile, err := os.Create(masterPath)
			if err != nil {
				log.Fatalf("Error creating master playlist file: %s", err)
			}
			defer masterFile.Close()

			if _, err := masterBuf.WriteTo(masterFile); err != nil {
				log.Fatalf("Error writing master playlist: %s", err)
			}
		}
	}
}

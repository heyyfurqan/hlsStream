package video

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/grafov/m3u8"
)

func CreateMasterPlaylist(outputDirectory string, resolutions []string, bitrates []string) {
	masterPl := m3u8.NewMasterPlaylist()

	for i, resolution := range resolutions {
		bitrate := bitrates[i]
		playlistPath := "playlist_" + resolution + ".m3u8"

		mediaPl, err := m3u8.NewMediaPlaylist(3, 3)
		if err != nil {
			panic(err)
		}

		params := m3u8.VariantParams{
			Bandwidth:  ResolutionToBandwidth(bitrate),
			Resolution: resolution,
		}

		masterPl.Append(playlistPath, mediaPl, params)
	}

	// Create the master playlist file in the output directory
	f, err := os.Create(filepath.Join(outputDirectory, "master.m3u8"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := masterPl.Encode()
	_, err = buf.WriteTo(f)
	if err != nil {
		panic(err)
	}
}

func ResolutionToBandwidth(bitrate string) uint32 {
	// Convert the bitrate to an integer
	bitrateInt, err := strconv.Atoi(bitrate)
	if err != nil {
		log.Fatalf("Invalid bitrate: %v", err)
	}

	// Calculate the bandwidth by adding 15% to the bitrate
	bandwidth := uint32(float64(bitrateInt) * 1.15)

	return bandwidth
}

# Simple http server to test hls stream

## 1- Choose media files (video and audio)
Select a random .mp4 video and audio file according to [Apple's standards](https://developer.apple.com/documentation/http-live-streaming/hls-authoring-specification-for-apple-devices#Video)

## 2- Project structure
Here is the file structure that this project uses:
```
.
├── audio
│   ├── audioProcessor.go
│   ├── segments/
│   └── storage/
│   └── player/
├── video
│   ├── playlist.go
│   └── segment.go
│   ├── segments/
│   └── storage/
│   └── player/
├── main.go
├── go.mod
└── go.sum
└── server.go
```
Place your media files in ```audio/storage/``` and ```video/storage/```.

## 3- Start the go server
Run the go server using:
```go run main.go server.go```
This not only starts the http server, but also segments the media into different resolutions, creates playlist files (*.m3u8) for them inside the folder titled as media's name as well as a master playlist (master.m3u8) that references these streams.

After the files have been segmented and the playlist files have been created, make sure to edit the hlsPlayers (```video/player/hlsStream.html``` && ```audio/player/hlsPlayer.html```) and update the ```src``` to the ```master.m3u8``` files, then go to [http://localhost:8080/](http://localhost:8080/), and open the HTML file. It will play the hls file.

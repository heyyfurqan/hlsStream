# Simple http server to test hls stream

## 1- Choose a video
Select a random .mp4 video according to [Apple's standards](https://developer.apple.com/documentation/http-live-streaming/hls-authoring-specification-for-apple-devices#Video)

## 2- Segmentation using FFmpeg
The following commands segment the video into 3 different resolutions, creates 3 different playlist files for them (stream.m3u8) as well as a master playlist (master.m3u8) that references these streams.

```
ffmpeg -i video.mp4 -filter_complex "[0:v]split=3[v1][v2][v3]; [v1]copy[v1out]; [v2]scale=w=1280:h=720[v2out];[v3]scale=w=640:h=360[v3out]" \
-map "[v1out]" -c:v:0 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:0 5M -maxrate:v:0 5M -minrate:v:0 5M -bufsize:v:0 10M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
-map "[v2out]" -c:v:1 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:1 3M -maxrate:v:1 3M -minrate:v:1 3M -bufsize:v:1 3M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
-map "[v3out]" -c:v:2 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:2 1M -maxrate:v:2 1M -minrate:v:2 1M -bufsize:v:2 1M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
-map a:0 -c:a:0 aac -b:a:0 96k -ac 2 \
-map a:0 -c:a:1 aac -b:a:1 96k -ac 2 \
-map a:0 -c:a:2 aac -b:a:2 48k -ac 2 \
-f hls \
-hls_time 2 \
-hls_playlist_type vod \
-hls_flags independent_segments \
-hls_segment_type mpegts \
-hls_segment_filename stream_%v/data%02d.ts \
-var_stream_map "v:0,a:0 v:1,a:1 v:2,a:2" \
-master_pl_name master.m3u8 \
stream_%v/stream.m3u8
```
## 3- Test on live server
Run the go server, then go to [http://localhost:8080/](http://localhost:8080/), and open the HTML file. It will play the hls file.

# gst-webrtc

```
sudo apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-libav
```

or

```
brew install gst-plugins-good gst-plugins-ugly pkg-config && export PKG_CONFIG_PATH="/usr/local/opt/libffi/lib/pkgconfig"
```

## ffmpeg

```
ffmpeg -re -r 30 -f lavfi -i testsrc -vf scale=1280:960 -vcodec libx264 -profile:v baseline -pix_fmt yuv420p -f flv rtmp://localhost/live/test
```

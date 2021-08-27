package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/pion/webrtc/v3"

	gst "github.com/muxable/gst-webrtc/gstreamer-src"
	"github.com/muxable/gst-webrtc/signal"
)

func demo(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func sdp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion1")

	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(audioTrack)
	if err != nil {
		panic(err)
	}

	firstVideoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion2")
	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(firstVideoTrack)

	secondVideoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion3")
	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(secondVideoTrack)
	if err != nil {
		panic(err)
	}

	offer := webrtc.SessionDescription{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	signal.Decode(buf.String(), &offer)

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	<-gatherComplete

	fmt.Fprintf(w, signal.Encode(*peerConnection.LocalDescription()))

	//	gst.CreatePipeline("opus", []*webrtc.TrackLocalStaticSample{audioTrack}, "rtmpsrc location=\"rtmp://localhost/live/test live=1\" ! flvdemux name=demux demux.audio ! queue ! decodebin").Start()
	gst.CreatePipeline("vp8", []*webrtc.TrackLocalStaticSample{firstVideoTrack, secondVideoTrack}, "rtmpsrc location=\"rtmp://localhost/live/test live=1\" ! flvdemux name=demux demux.video").Start()
}

func main() {
	http.HandleFunc("/", demo)
	http.HandleFunc("/sdp", sdp)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

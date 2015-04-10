#!/bin/bash

cd /root/horus-v2/bin

# set the autofocus on
v4l2-ctl -d /dev/video1 -c focus_auto=1

# Next line not necessary if you are using my -F option on capture
v4l2-ctl --set-fmt-video=width=1920,height=1080,pixelformat=1

# Pipe the output of capture into avconv/ffmpeg
# capture "-F"   My H264 passthrough mode
#         "-o"   Output the video (to be passed to avconv via pipe)
#         "-c0"  Capture 0 frames, which means infinite frames in my program
# avconv "-re" read input at the native frame rate
#        "-i -"  Take the input from the pipe
#        "-vcodec copy" Do not transcode the video
#        "-f rtp rtp://192.168.1.2:1234/" Force rtp to output to address of my PC on port 1234
./capture -d /dev/video1 -F -o -c0|avconv -re -i - -vcodec copy -f flv rtmp://publish-sfo1.cine.io/live/XJRl3Bsq?group40

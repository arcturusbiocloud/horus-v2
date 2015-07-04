#!/bin/bash

cd /home/pi/horus-v2/bin

# Next line not necessary if you are using my -F option on capture
v4l2-ctl --set-fmt-video=width=1920,height=1080,pixelformat=1

# set the C920 settings to get the picture
# http://askubuntu.com/questions/211971/v4l2-ctl-exposure-auto-setting-fails
# http://stackoverflow.com/questions/13407859/is-there-a-way-to-control-a-webcam-focus-in-pygame

if [ $1 = "UV" ]; then
  # write 1 to the arduino serial port to turn on the uv light
  # write 0 to the arduino serial port to turn off the uv light
  echo "Taking a picture with UV light"
  v4l2-ctl -d /dev/video0 -c focus_auto=0
  v4l2-ctl -d /dev/video0 -c focus_absolute=70
  v4l2-ctl --device=/dev/video0 --set-ctrl=exposure_auto=1
  v4l2-ctl --device=/dev/video0 --set-ctrl=exposure_absolute=2500
  ./boneCV 0
fi

if [ $1 = "WHITE" ]; then
 # write B to the arduino serial port to turn on the white light
 # write C to the arduino serial port to turn off the white light 
 echo "Taking a picture with WHITE light"
 v4l2-ctl -d /dev/video3 -c focus_auto=0
 v4l2-ctl -d /dev/video3 -c focus_absolute=85
 v4l2-ctl --device=/dev/video3 --set-ctrl=exposure_auto=1
 v4l2-ctl --device=/dev/video3 --set-ctrl=exposure_absolute=0
 ./boneCV 3
fi


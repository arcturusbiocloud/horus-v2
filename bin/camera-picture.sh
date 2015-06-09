#!/bin/bash

cd ~/horus-v2/bin

# Next line not necessary if you are using my -F option on capture
v4l2-ctl --set-fmt-video=width=1920,height=1080,pixelformat=1


# set the C920 settings to get the picture
# http://askubuntu.com/questions/211971/v4l2-ctl-exposure-auto-setting-fails
# http://stackoverflow.com/questions/13407859/is-there-a-way-to-control-a-webcam-focus-in-pygame
v4l2-ctl -d /dev/video0 -c focus_auto=0
v4l2-ctl -d /dev/video0 -c focus_absolute=50

./boneCV 0
g++ -O2 `pkg-config --cflags --libs opencv` boneCV.cpp -o boneCV
gcc -O2 -Wall `pkg-config --cflags --libs libv4l2` capture.c -o capture
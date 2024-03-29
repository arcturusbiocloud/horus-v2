/* boneCV.cpp
 *
 * Copyright Derek Molloy, School of Electronic Engineering, Dublin City University
 * www.derekmolloy.ie
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted
 * provided that source code redistributions retain this notice.
 *
 * This software is provided AS IS and it comes with no warranties of any type. 
 */

/*
 * sudo apt-get install libopencv-dev
 * build with: g++ -O2 `pkg-config --cflags --libs opencv` boneCV.cpp -o boneCV
 *
 */

#include<iostream>
#include<opencv2/opencv.hpp>
using namespace std;
using namespace cv;

int main(int argc, char ** argv)
{
    int video = 0;
    
    // parsing parameter
    if (argc > 1) video = atoi(argv[1]);
    
    VideoCapture capture(video);
    capture.set(CV_CAP_PROP_FRAME_WIDTH,1920);
    capture.set(CV_CAP_PROP_FRAME_HEIGHT,1080);
    if(!capture.isOpened()){
	    cout << "Failed to connect to the camera." << endl;
    }
    // Mat frame, edges;
    Mat frame;
    capture >> frame;
    if(frame.empty()){
		cout << "Failed to capture an image" << endl;
		return -1;
    }
    // cvtColor(frame, edges, CV_BGR2GRAY);
    // Canny(edges, edges, 0, 30, 3);
    // imwrite("edges.png", edges);
    imwrite("capture.png", frame);
    return 0;
}

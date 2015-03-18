Horus-v2
=====

This project is responsible for manage the execution of the scripts on the Arcturus BioCloud "biocomputer". It should evolve to something that can control all the hardware through the serial and USB ports of the BBB.

## Install dependencies
    go get github.com/go-martini/martini
    go get github.com/codegangsta/martini-contrib/render

## How to start
    go run main.go
    
## How to test
    # POST /api/project/script_call
    curl -X POST localhost:3000/api/project/experiment.sh%201
    {"script_call":"experiment.sh 1","status":"running"}

    # GET /api/project/script_call
    curl -X GET localhost:3000/api/project/experiment.sh%201
    {"script_call":"experiment.sh 1","status":"running"}
    
    # GET /api/online
    {"status":"true"}

## Feature Roadmap

  - [ ] rest call to execute and check bash scripts status
  - [ ] rest call to check the status of the hardware
  - [ ] basic authentication
  - [ ] auto init BBB script
  - [ ] ARM cross compilation
  - [ ] rest call to start video streaming
  - [ ] storage device interface (OpenPCR)
  - [ ] serial port interface (labcontrol)
  - [ ] robot syntax sugar
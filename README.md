Horus-v2
=====

This project is responsible for manage the execution of the scripts on the Arcturus BioCloud "biocomputer". It should evolve to something that can control all the hardware through the serial and USB ports of the BBB.

## Install dependencies
    go get github.com/go-martini/martini
    go get github.com/codegangsta/martini-contrib/render

## How to start
    go run main.go
    
## How to cross compile to BBB
You must first configure your Go environment for arm linux cross compiling

    $ cd $GOROOT/src
    $ GOOS=linux GOARCH=arm ./make.bash --no-clean
    
Then compile your Gobot program with

    $ GOARM=7 GOARCH=arm GOOS=linux go build
    
## How to install
    
If you are running the official Angstrom or Debian linux through the usb->ethernet connection, you can simply upload your program and install it with

    $ scp horus-v2 root@192.168.7.2:/usr/bin/
    $ scp horus.service root@192.168.7.2:/lib/systemd/
    $ ssh root@192.168.7.2
    $ cd /etc/systemd/system/
    $ ln /lib/systemd/horus.service horus.service
    $ systemctl daemon-reload
    $ systemctl start horus.service
    $ systemctl enable horus.service
    
## How to clean the BBB

* <http://kacangbawang.com/beagleboneblack-revc-debloat-part-1/>
    
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

  - [x] ARM cross compilation
  - [x] auto init BBB script
  - [ ] rest call to execute and check bash scripts status
  - [ ] rest call to check the status of the hardware
  - [ ] basic authentication
  - [ ] rest call to start video streaming
  - [ ] storage device interface (OpenPCR)
  - [ ] serial port interface (labcontrol)
  - [ ] robot syntax sugar
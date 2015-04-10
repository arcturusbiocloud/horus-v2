Horus-v2
=====

This project is responsible for manage the execution of the scripts on the Arcturus BioCloud "biocomputer". It should evolve to something that can control all the hardware through the serial and USB ports of the BBB.

## Install dependencies
    go get github.com/go-martini/martini
    go get github.com/codegangsta/martini-contrib/render
    go get github.com/martini-contrib/auth
    go get github.com/tarm/serial
    go get github.com/mitchellh/go-ps

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

    $ scp horus-v2 root@10.1.10.111:/usr/bin/
    $ scp horus.service root@10.1.10.111:/lib/systemd/
    $ ssh root@10.1.10.111
    $ cd /etc/systemd/system/
    $ ln /lib/systemd/horus.service horus.service
    $ systemctl daemon-reload
    $ systemctl start horus.service
    $ systemctl enable horus.service
    
## How to clean the BBB

* <http://kacangbawang.com/beagleboneblack-revc-debloat-part-1/>
    
## How to test    
    # GET /api/online
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/online
    {"status":"Humidity: 35.40% Temperature: 29.90C"}
    # or
    {"status":"error", "error: ... "}
    
    # turn camera 0 streaming on
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/camera_streaming/on
    {"status":"streaming"}
    # or
    {"status":"error", "error: ... "}
    
    # turn camera 0 streaming off
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/camera_streaming/off
    {"status": "streaming stopped"}
    # or
    {"status":"error", "error: ..."}
    
    # get a picture from camera 0 at a specific slot and upload to a specific project at arcturus.io
    # /api/take_picture/:project_id/:petri_dish_slot/uv_on|uv_off/light_on|light_off
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/take_picture/2/5/uv_on/light_off
    {"status":"Taking picture for the project 2 at the petri dish slot 5"}
    # or
    {"status":"error", "error":"Machine already ocuppied by another process."}
        
## Feature Roadmap

  - [x] ARM cross compilation
  - [x] auto init BBB script
  - [x] serial port interface transilluminator
  - [x] serial port interface incubator
  - [x] rest call to start video streaming
  - [x] rest call to take a picture from a specific slot
  - [x] same response format to all calls
  - [x] syntax sugar to serial devices
  - [x] basic authentication
  - [x] rest call to check the status of the hardware
  - [ ] rest call to make a experiment
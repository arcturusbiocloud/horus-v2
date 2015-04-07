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
    curl --user arcturus:huxnGrbNfQFR -X POST localhost:3000/api/project/sleep%2060
    {"status":"running", "pid":14298,"script_call":"sleep 60"}

    # GET /api/project/pid
    curl --user arcturus:huxnGrbNfQFR -X GET localhost:3000/api/project/14224
    {"pid":14298,"status":"alive"}
    # or
    {"status":"error", "err":"os: process already finished","pid":14298}
    
    # GET /api/online
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/online
    {"status":"true"}
    # or
    {"status":"error" ...}
    
    # turn UV light OFF
    $ curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/uv_light/off
    {"status":"uv light turned off"}
    # or
    {"status":"error" ...}
    
    # turn UV light ON
    $ curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/uv_light/on
    {"status":"uv light turned on"}
    # or
    {"status":"error" ...}
    
    # get the temperature and humidity stats
    $ curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/incubator/stats
    {"status":"Humidity: 37.20 %\tTemperature: 30.00 *C"}
    # or
    {"status":"error" ...}
    
    # turn camera 0 streaming on
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/camera_streaming/on
    {"status":"streaming"}
    # or
    {"status":"error" ...}
    
    # turn camera 0 streaming off
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/camera_streaming/off
    {"status": "streaming stopped"}
    # or
    {"status":"error" ...}
    
    # get a picture from camera 1 at a specific slot
    curl --user arcturus:huxnGrbNfQFR -X GET 10.1.10.111:3000/api/camera_picture/1-11
    # the response is a png file or a 500 internal server error. it stops the camera 0 streaming
        
## Feature Roadmap

  - [x] ARM cross compilation
  - [x] auto init BBB script
  - [x] serial port interface transilluminator
  - [x] serial port interface incubator
  - [x] rest call to execute and check bash scripts status
  - [x] rest call to start video streaming
  - [x] rest call to take a picture from a specific slot
  - [x] same response format to all calls
  - [x] syntax sugar to serial devices
  - [x] basic authentication
  - [ ] rest call to check the status of the hardware
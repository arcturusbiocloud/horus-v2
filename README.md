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
    
## How to install on Rpi2

If you are running the Raspbian

    $ sudo cp init-scripts/horus-v2.sh /etc/init.d/
    $ sudo update-rc.d horus-v2.sh defaults    
    
## How to clean the BBB

* <http://kacangbawang.com/beagleboneblack-revc-debloat-part-1/>
    
## How to test    
    # check the Arduino controller returning the temperature and humidity
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.42:3000/api/online
    {"status":"Humidity: 35.40% Temperature: 29.90C"}
    # or
    {"status":"error", "error: ... "}
    
    # initialize the OpenPCR
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.42:3000/api/init_pcr
    {"status":"OpenPCR initialized"}
    
    # zero the Modular Science robot
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.42:3000/api/zero_machine
    {"status": "Modular science robot zeroed"}
    
    # turn camera 1 streaming on
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.42:3000/api/camera_streaming/on
    {"status":"streaming"}
    # or
    {"status":"error", "error: ... "}
    
    # turn camera 1 streaming off
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.42:3000/api/camera_streaming/off
    {"status": "streaming stopped"}
    # or
    {"status":"error", "error: ..."}
    
    # get a picture from camera 0 at a specific slot and upload to a specific project at arcturus.io
    # /api/take_picture/:project_id/:petri_dish_slot/uv_on|uv_off/light_on|light_off
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.42:3000/api/take_picture/2/5/uv_on/light_off
    {"status":"Taking picture for the project 2 at the petri dish slot 5"}
    # or
    {"status":"error", "error":"Machine already ocuppied by another process."}
     
    # run a experiment
    curl --user arcturus:huxnGrbNfQFR --data "project_id=4&slot=5&genetic_parts={'anchor': 'clhor', 'promoter': 'very_strong', 'rbs': 'very_strong', 'gene': 'gfp', 'terminator': 'ter', 'cap': 'high'}" -v -X POST 192.168.7.42:3000/api/run_experiment
    {"status": "Running experiment for the project 1 at the petri dish slot 5 or 6 with the genetic parts ..."}
    # or
    {"status": "error", "error": "Machine already ocuppied by another process."}
    # or
    {"status": "error", "error": "Petri dish slot out of range."}
    
    # turn on the centrifuge
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.41:3000/api/centrifuge/on
    {"status":"5\r\n"}
    # or
    {"status":"error", "error: ... "}
    
    # turn off the centrifuge
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.41:3000/api/centrifuge/off
    {"status":"6\r\n"}
    # or
    {"status":"error", "error: ... "}
    
    # turn on the shaker
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.41:3000/api/shaker/on
    {"status":"7\r\n"}
    # or
    {"status":"error", "error: ... "}
    
    # turn off the shaker
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.41:3000/api/shaker/off
    {"status":"8\r\n"}
    # or
    {"status":"error", "error: ... "}
    
    # turn on the gel
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.41:3000/api/gel/on
    {"status":"9\r\n"}
    # or
    {"status":"error", "error: ... "}
    
    # turn off the gel
    curl --user arcturus:huxnGrbNfQFR -X GET 192.168.7.41:3000/api/gel/off
    {"status":"A\r\n"}
    # or
    {"status":"error", "error: ... "}
    
    # get petri dish picture on the uv light chamber
    http://arcturus:huxnGrbNfQFR@horus01.arcturus.io:3001/api/camera_picture_petri_dish/uv
    
    # get petri dish picture on the white light chamber
    http://arcturus:huxnGrbNfQFR@horus01.arcturus.io:3001/api/camera_picture_petri_dish/white
    
     
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
  - [x] rest call to make a experiment
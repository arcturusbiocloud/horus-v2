package main

import (
  "github.com/go-martini/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/tarm/serial"
  "log"
  "time"
  "bufio"
)

func main() {
  m := martini.Classic()
  m.Use(render.Renderer())
  
  // init the arduino serial port
  c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600, ReadTimeout: time.Millisecond * 2000}
  s, _ := serial.OpenPort(c)
  log.Printf("The USB port the Arduino service will use is /dev/ttyACM0")
  
  // When connecting to an older revision Arduino, you need to wait
	// a little while it resets.
	time.Sleep(1 * time.Second)
	  
  // execute a script and return the status
  m.Post("/api/project/:script_call", func(r render.Render, params martini.Params) {
    script_call := params["script_call"]
    // run a script
    // ...
    // status will be running, not_found
    r.JSON(200, map[string]interface{}{"script_call": script_call, "status": "running"})
  })
  
  // get a status of a script previously executed
  m.Get("/api/project/:script_call", func(r render.Render, params martini.Params) {
    script_call := params["script_call"]
    // list process with the name script_call and check if it's running
    // if running
    r.JSON(200, map[string]interface{}{"script_call": script_call, "status": "running"})
    // if not
    // ...
  })

  // check all the machine hardware
  m.Get("/api/online", func(r render.Render) {
    // test stack hardware
    // ...
    r.JSON(200, map[string]interface{}{"status": "true"})
  })
  
  // send commands to the arduino serial port
  m.Get("/api/serial/:buffer", func(r render.Render, params martini.Params) {
    buf := params["buffer"]
        
    // check the serial port
    if s == nil {
      r.JSON(200, map[string]interface{}{"error": "serial port not connected"})
      return
    }
    
    // get serial
    _, err := s.Write([]byte(buf))
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"error": err})
    } else {
      bio := bufio.NewReader(s)
      buf, err := bio.ReadString('\n')
      if err != nil {
        r.JSON(200, map[string]interface{}{"error": err})
      } else {
        r.JSON(200, map[string]interface{}{"status": string(buf)})
      }
    }
  })
    
  m.Run()
}

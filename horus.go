package main

import (
  "github.com/go-martini/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/tarm/serial"
  "github.com/mitchellh/go-ps"
  "time"
  "bufio"
  "log"
  "strings"
  "os"
  "io"
  "strconv"
  "syscall"
  "os/exec"
  "net/http"
)

func main() {
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Use(martini.Static("/root/horus-v2/streaming"))
    
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
    pid, err := exe_cmd(script_call)
    if err != nil {
      r.JSON(200, map[string]interface{}{"script_call": script_call, "status": "error", "err" : err})
    } else {
      // status will be running, not_found
      r.JSON(200, map[string]interface{}{"script_call": script_call, "status": "running", "pid": pid})
    }
  })
  
  // get a status of a script previously executed
  m.Get("/api/project/:pid", func(r render.Render, params martini.Params) {
    pid := params["pid"]
    
    iPid, _ := strconv.Atoi(pid)
    p, _ := os.FindProcess(iPid)
    err := p.Signal(syscall.Signal(0))
    if err != nil {
      r.JSON(200, map[string]interface{}{"pid": iPid, "status": "error", "err": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"pid": iPid, "status": "alive"})
    }
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
  
  m.Get("/api/camera_streaming/on", func(r render.Render) {    
    // kill any previous streaming
    turnoff_streaming()
    
    _, err := exe_cmd("/root/horus-v2/bin/camera-streaming.sh")
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": "streaming"})
    }
  })
  
  m.Get("/api/camera_streaming/off", func(r render.Render) {        
    turnoff_streaming()
    r.JSON(200, map[string]interface{}{"status": "streaming stopped"})
  })
  
  m.Get("/api/camera_picture", func(res http.ResponseWriter, req *http.Request) {
    turnoff_streaming()
    
    // remove files
    os.Remove("/root/horus-v2/bin/capture.png")
    os.Remove("/root/horus-v2/bin/edges.png")
    
    // take picture
    proc := exec.Command("v4l2-ctl", "-d", "/dev/video1", "-c", "focus_auto=0")
    proc.Run()
    proc = exec.Command("v4l2-ctl", "--set-fmt-video=width=1920,height=1080,pixelformat=1")
    proc.Run()
    proc = exec.Command("/root/horus-v2/bin/boneCV")
    err := proc.Run()
    if err != nil {
      log.Printf("err= %s", err.Error())
      res.WriteHeader(500)
    }

    // serving image
    res.Header().Set("Content-Type", "image/png")
    f, err := os.Open("/root/horus-v2/bin/capture.png")
    if err != nil {
      log.Printf("err= %s", err.Error())
      res.WriteHeader(500)
    }
    defer f.Close()
    io.Copy(res, f)
  })
    
  m.Run()
}

func turnoff_streaming() {
  p, _ := ps.Processes()
  for _, p1 := range p {
    if p1.Executable() == "camera-streamin" || p1.Executable() == "capture" ||  p1.Executable() == "avconv" {
      proc, _ := os.FindProcess(p1.Pid())
      proc.Kill()
    }
  }
}

func exe_cmd(cmd string) (int,error) {
  parts := strings.Fields(cmd)
  head := parts[0]
  parts = parts[1:len(parts)]
  
  out := exec.Command(head, parts...)
  err := out.Start()
  if err != nil {
    return 0, err
  }
  go out.Wait()
  return out.Process.Pid, nil
}

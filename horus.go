package main

import (
  "github.com/go-martini/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/martini-contrib/auth"
  "github.com/tarm/serial"
  "github.com/mitchellh/go-ps"
  "os/exec"
  "net/http"
  "time"
  "bufio"
  "log"
  "strings"
  "os"
  "io"
  "strconv"
  "syscall"
  "errors"
)

// serial port handle to communicate with the Arduino
var s *serial.Port
// handle to check if the machine is already executing some script 
var running = false

func main() {
  log.Printf("Horus-v2 bio server controller")
  
  // config martini handlers
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Use(martini.Static("/root/horus-v2/streaming"))

  // authenticate every request
  m.Use(auth.BasicFunc(func(username, password string) bool {
    return username == "arcturus" && password == "huxnGrbNfQFR"
  }))

  // check in every call if the machine is already executing some script
  m.Use(func(r render.Render) {
    if (running) {
      r.JSON(200, map[string]interface{}{"status": "error", "error": "Machine already ocuppied by another process."})
    }
  })
      
  // init the arduino serial port
  c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600, ReadTimeout: time.Millisecond * 2000}
  s, _ = serial.OpenPort(c)
  
  // when connecting to an older revision Arduino, you need to wait a little while it resets.
	time.Sleep(1 * time.Second)
		  
  // execute a script and return the status
  m.Post("/api/project/:script_call", func(r render.Render, params martini.Params) {
    script_call := params["script_call"]
    // run a script
    pid, err := exe_cmd(script_call)
    if err != nil {
      r.JSON(200, map[string]interface{}{"script_call": script_call, "status": "error", "error" : err.Error()})
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
      r.JSON(200, map[string]interface{}{"pid": iPid, "status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"pid": iPid, "status": "alive"})
    }
  })

  // check all the machine hardware
  m.Get("/api/online", func(r render.Render) {
    // test stack hardware
    // ... send a ping to the arduino controller and check the robot online
    r.JSON(200, map[string]interface{}{"status": "true"})
  })
  
  // switch on and off the incubator and the tent light
  m.Get("/api/light/:set", func(r render.Render, params martini.Params) {
    var response_buf string
    var err error
    
    if params["set"] == "on" {
      response_buf, err = turn_on_light()
    } else if params["set"] == "off" {
      response_buf, err = turn_off_light()
    } else {
      r.JSON(200, map[string]interface{}{"status": "error", "error": "Parameter not supported."})
      return
    }
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      if string(response_buf) == "3\r\n" || string(response_buf) == "4\r\n" {
        r.JSON(200, map[string]interface{}{"status": "light switched"})
      } else {
        r.JSON(200, map[string]interface{}{"status": "error", "error": "unexpected response " + string(response_buf)})
      }
    } 
  })  
  
  // turn on the uv light at same time turning off the incubator and the tent light
  m.Get("/api/uv_light/on", func(r render.Render, params martini.Params) {
    response_buf, err := turn_on_uv_light()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      if string(response_buf) != "1\r\n" {
        r.JSON(200, map[string]interface{}{"status": "error", "error": "unexpected response " + string(response_buf)})
      } else {
        r.JSON(200, map[string]interface{}{"status": "uv light turned on"})
      }
    } 
  })

  // turn off the uv light at same time turning on the incubator and the tent light
  m.Get("/api/uv_light/off", func(r render.Render, params martini.Params) {
    response_buf, err := turn_off_uv_light()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      if string(response_buf) != "0\r\n" {
        r.JSON(200, map[string]interface{}{"status": "error", "error": "unexpected response " + string(response_buf)})
      } else {
        r.JSON(200, map[string]interface{}{"status": "uv light turned off"})
      }
    }
  })

  // get the humidity and the temperature from the incubator
  m.Get("/api/incubator/stats", func(r render.Render, params martini.Params) {
    response_buf, err := get_incubator_stats()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": strings.Replace(string(response_buf), "\r\n", "", 1)})
    }
  })
  
  // turn on the camera 0 live streaming
  m.Get("/api/camera_streaming/on", func(r render.Render) {    
    // kill any previous streaming
    turn_off_streaming()
    
    _, err := exe_cmd("/root/horus-v2/bin/camera-streaming.sh")
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": "streaming"})
    }
  })
  
  // turn off the camera 0 live streaming
  m.Get("/api/camera_streaming/off", func(r render.Render) {        
    turn_off_streaming()
    r.JSON(200, map[string]interface{}{"status": "streaming stopped"})
  })
  
  // take a picture using the camera 1
  m.Get("/api/camera_picture/:slot/:uv", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
    // we cannot use two cameras at the same time
    turn_off_streaming()
    
    // remove files
    os.Remove("/root/horus-v2/bin/capture.png")
    
    // define slot
    slot := params["slot"]
    iSlot, _ := strconv.Atoi(slot)
    if iSlot <= 0 || iSlot >=12 {
      res.WriteHeader(500)
      return
    }
        
    // run scripts to open the oven, positioning on the grid, open the petri dish
    // python /root/labcontrol/labcontrol.py -S 1 -v -w /root/labcontrol -s openOven_openPetriDish_putCamera.py
    running = true
    proc := exec.Command("python", "/root/labcontrol/labcontrol.py", "-S", slot, "-v", "-w", "/root/labcontrol", "-s", "openOven_openPetriDish_putCamera.py")
    proc.Run()
    
    // turn on the UV light
    uv := params["uv"]
    if uv == "uv_on" {
      turn_on_uv_light()
    } else {
      turn_off_light()
    }
    
    // take picture
    // http://askubuntu.com/questions/211971/v4l2-ctl-exposure-auto-setting-fails
    // http://stackoverflow.com/questions/13407859/is-there-a-way-to-control-a-webcam-focus-in-pygame
    proc = exec.Command("v4l2-ctl", "--set-fmt-video=width=1920,height=1080,pixelformat=1")
    proc.Run()
    proc = exec.Command("v4l2-ctl", "-d", "/dev/video1", "-c", "focus_auto=0")
    proc.Run()
    proc = exec.Command("v4l2-ctl", "-d", "/dev/video1", "-c", "focus_absolute=50")
    proc.Run()
    proc = exec.Command("/root/horus-v2/bin/boneCV")
    err := proc.Run()
    if err != nil {
      log.Printf("err= %s", err.Error())
      res.WriteHeader(500)
      running = false
      return
    }
    
    // turn on the UV light
    if uv == "uv_on" {
      turn_off_uv_light()
    } else {
      turn_on_light()
    }
        
    // close the  petri dish, turn off the UV light, close the oven, go home
    // python /root/labcontrol/labcontrol.py -S 1 -v -w /root/labcontrol -s closePetriDish_closeOven_goHome.py
    proc = exec.Command("python", "/root/labcontrol/labcontrol.py", "-S", slot, "-v", "-w", "/root/labcontrol", "-s", "closePetriDish_closeOven_goHome.py")
    proc.Run()
    
    f, err := os.Open("/root/horus-v2/bin/capture.png")
    if err != nil {
      log.Printf("err= %s", err.Error())
      res.WriteHeader(500)
      running = false
      return
    }
    defer f.Close()
    
    // serving image
    res.Header().Set("Content-Type", "image/png")
    io.Copy(res, f)
    running = false
  })
    
  m.Run()
}

func turn_off_streaming() {
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

func turn_on_uv_light() (string, error) {
  return serial_cmd("1")
}

func turn_off_uv_light() (string, error) {
  return serial_cmd("0")
}

func get_incubator_stats() (string, error) {
  return serial_cmd("2")
}

func turn_off_light() (string, error) {
  return serial_cmd("3")
}

func turn_on_light() (string, error) {
  return serial_cmd("4")
}

func serial_cmd(cmd string) (string, error) {
  // check the serial port
  if s == nil {
    return "", errors.New("serial port not connected")
  }
  
  // get serial
  _, err := s.Write([]byte(cmd))
  
  if err != nil {
    return "", err
  } else {
    bio := bufio.NewReader(s)
    buf, err := bio.ReadString('\n')
    if err != nil {
      return "", err
    } else {
      return string(buf), nil
    }
  }
}

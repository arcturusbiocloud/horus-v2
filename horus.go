package main

import (
  "github.com/go-martini/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/martini-contrib/auth"
  "github.com/tarm/serial"
  "github.com/mitchellh/go-ps"
  "net/http"
  "os/exec"
  "time"
  "bufio"
  "log"
  "strings"
  "os"
  "io"
  "strconv"
  "fmt"
  "errors"
  "path/filepath"
)

// serial port handle to communicate with the Arduino
var s *serial.Port

// handle to check if the machine is already executing some script 
var running = false
var dir = ""

func main() {
  // get execution path
  dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
  log.Printf("Horus-v2 bio server controller running at %s", dir)
  
  // config martini handlers
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Use(martini.Static(dir + "/streaming"))

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
  c := &serial.Config{Name: "/dev/ttyACM1", Baud: 9600, ReadTimeout: time.Millisecond * 2000}
  s, _ = serial.OpenPort(c)
  
  // when connecting to an older revision Arduino, you need to wait a little while it resets.
  time.Sleep(1 * time.Second)
		  
  // check all the machine hardware
  m.Get("/api/online", func(r render.Render) {
    running = true
    // test stack hardware
    response_buf, err := get_incubator_stats()
    running = false
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": strings.Replace(string(response_buf), "\r\n", "", 1)})
    }
  })
  
  // zero the machine. WARNING: the robot should be at the proper position
  m.Get("/api/zero_machine", func(r render.Render) {    
    running = true
    proc := exec.Command("python", filepath.Dir(dir) + "/labcontrol/labcontrol.py", "-v", "-w", filepath.Dir(dir) + "/labcontrol", "-s", "zero.py")
    proc.Run()
    running = false
    
    r.JSON(200, map[string]interface{}{"status": "Modular science robot zeroed"})
  })
  
  // zero the machine. WARNING: the robot should be at the proper position
  m.Get("/api/init_pcr", func(r render.Render) {    
    proc := exec.Command("bash", filepath.Dir(dir) + "/OpenPyCR/simple-pcr-run.sh")
    proc.Run()
    
    r.JSON(200, map[string]interface{}{"status": "OpenPCR initialized"})
  })
  
  // turn on uv light
  m.Get("/api/uv_light/on", func(r render.Render) {    
    
    // turn on uv_light
    buf, err := turn_on_uv_light()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
  
  // turn off uv light
  m.Get("/api/uv_light/off", func(r render.Render) {    
    
    // turn off uv_light
    buf, err := turn_off_uv_light()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
  
  // turn on centrifuge
  m.Get("/api/centrifuge/on", func(r render.Render) {    
    
    // turn on centrifuge
    buf, err := turn_on_centrifuge()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
  
  // turn off centrifuge
  m.Get("/api/centrifuge/off", func(r render.Render) {    
    
    // turn off centrifuge
    buf, err := turn_off_centrifuge()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
  
  // turn on shaker
  m.Get("/api/shaker/on", func(r render.Render) {    
    
    // turn on shaker
    buf, err := turn_on_shaker()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
  
  // turn off shaker
  m.Get("/api/shaker/off", func(r render.Render) {    
    
    // turn off shaker
    buf, err := turn_off_shaker()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
  
  // turn on gel
  m.Get("/api/gel/on", func(r render.Render) {    
    
    // turn on gel
    buf, err := turn_on_gel()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
  
  // turn off gel
  m.Get("/api/gel/off", func(r render.Render) {    
    
    // turn off gel
    buf, err := turn_off_gel()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": buf})
    }
  })
    
  // turn on the camera 1 live streaming
  m.Get("/api/camera_streaming/on", func(r render.Render) {    
    // kill any previous streaming
    turn_off_streaming()
    
    // turn on streaming
    _, err := turn_on_streaming()
    
    if err != nil {
      r.JSON(200, map[string]interface{}{"status": "error", "error": err.Error()})
    } else {
      r.JSON(200, map[string]interface{}{"status": "streaming"})
    }
  })
    
  // turn off the camera 1 live streaming
  m.Get("/api/camera_streaming/off", func(r render.Render) {        
    turn_off_streaming()
    r.JSON(200, map[string]interface{}{"status": "streaming stopped"})
  })
  
  m.Get("/api/camera_picture_petri_dish/uv", func(res http.ResponseWriter, req *http.Request) {
    running = true
    
    turn_on_uv_light()
    time.Sleep(2 * time.Second)
    camera_picture_petri_dish("UV")
    turn_off_uv_light()
    
    running = false
    
    // serving image
    res.Header().Set("Content-Type", "image/png")
    f, err := os.Open(filepath.Dir(dir) + "/horus-v2/bin/capture.png")
    if err != nil {
      log.Printf("err= %s", err.Error())
      res.WriteHeader(500)
    }
    defer f.Close()
    io.Copy(res, f)
    
  })
  
  m.Get("/api/camera_picture_petri_dish/white", func(res http.ResponseWriter, req *http.Request) {
    running = true
    
    turn_on_white_light()
    time.Sleep(2 * time.Second)
    camera_picture_petri_dish("WHITE")
    turn_off_white_light()
    
    running = false
    
    // serving image
    res.Header().Set("Content-Type", "image/png")
    f, err := os.Open(filepath.Dir(dir) + "/horus-v2/bin/capture.png")
    if err != nil {
      log.Printf("err= %s", err.Error())
      res.WriteHeader(500)
    }
    defer f.Close()
    io.Copy(res, f)
  })
    
  // take a picture using the camera 0
  m.Get("/api/take_picture/:project_id/:slot/:uv_light/:light", func(r render.Render, params martini.Params) {        
    // we cannot use two cameras at the same time
    turn_off_streaming()
    
    // get project_id     
    project_id, _ := strconv.Atoi(params["project_id"])
    
    // get petri dish slot
    slot := params["slot"]
    iSlot, _ := strconv.Atoi(slot)
    if iSlot <= 0 || iSlot >=12 {
      r.JSON(200, map[string]interface{}{"status": "error", "error": "Petri dish slot out of range."})
      return
    }
    
    // get uv_light parameter
    var uv_light bool
    if params["uv_light"] == "uv_on" {
      uv_light = true
    } else {
      uv_light = false
    }
    
    // get light parameter
    var light bool
    if params["light"] == "light_on" {
      light = true
    } else {
      light = false
    }
    
    go camera_picture(project_id, slot, uv_light, light)
    
    r.JSON(200, map[string]interface{}{"status": fmt.Sprintf("Taking picture for the project %d at the petri dish slot %s", project_id, slot)})
  })
  
  m.Post("/api/run_experiment", func(req *http.Request, r render.Render, params martini.Params) {
    // get project_id   
    project_id := req.FormValue("project_id")
    
    // get petri dish slot
    slot := req.FormValue("slot")
    iSlot, _ := strconv.Atoi(slot)
    if iSlot <= 0 || iSlot >=12 {
      r.JSON(200, map[string]interface{}{"status": "error", "error": "Petri dish slot out of range."})
      return
    }
    
    // get genetic parts
    genetic_parts := req.FormValue("genetic_parts")
                
    go run_experiment(project_id, slot, genetic_parts)
    
    r.JSON(200, map[string]interface{}{"status": fmt.Sprintf("Running experiment for the project %s at the petri dish slot %s with the genetic parts %s", project_id, slot, genetic_parts)})
  })
  
  m.Post("/api/run_virtual_experiment", func(req *http.Request, r render.Render, params martini.Params) {
    // get project_id   
    project_id := req.FormValue("project_id")
    
    // get petri dish slot
    slot := req.FormValue("slot")
    iSlot, _ := strconv.Atoi(slot)
    if iSlot <= 0 || iSlot >=12 {
      r.JSON(200, map[string]interface{}{"status": "error", "error": "Petri dish slot out of range."})
      return
    }
    
    // get genetic parts
    genetic_parts := req.FormValue("genetic_parts")
                
    go run_virtual_experiment(project_id, slot, genetic_parts)
    
    r.JSON(200, map[string]interface{}{"status": fmt.Sprintf("Running experiment for the project %s at the petri dish slot %s with the genetic parts %s", project_id, slot, genetic_parts)})    
  })
      
  m.Run()
}

func run_virtual_experiment(project_id string, slot string, genetic_parts string) (error) {
  // set the running state
  running = true
  
  // turn on virtual streaming
  go func() {
    exe_cmd("avconv -i " + filepath.Dir(dir) + "/horus-v2/streaming/XJRl3Bsq.20150411T022627.mp4 -vcodec copy -f flv rtmp://publish-sfo1.cine.io/live/XJRl3Bsq?group40")
  }()
  
  // kill the streaming after 5 minutes
  go func() {
    time.Sleep(300 * time.Second)
    
    // get the latest video processed by cine.io. It takes some time to encode the video
    // by now we are not getting the final video but adding a previous recorded
    // update the project with the final video
    proc := exec.Command("curl",
                          "--insecure",
                          "-X", "PUT", fmt.Sprintf("https://www.arcturus.io/api/projects/%s?access_token=55d28fc5783172b90fea425a2312b95a&recording_file_name=XJRl3Bsq.20150411T022627.mp4", project_id))
    _, err := proc.CombinedOutput()
    if err != nil {
      fmt.Printf("run_experiment() project_id=%d err=%s\n", project_id, err.Error())
    }                                                     
  }()
  
  // send the assembly update to arcturus.io project timeline
  proc := exec.Command("curl", 
                       "--insecure", 
                       "-X", "POST", fmt.Sprintf("https://www.arcturus.io/api/projects/%s/activities?access_token=55d28fc5783172b90fea425a2312b95a&key=1", project_id))
  _, err := proc.CombinedOutput()
  if err != nil {
    fmt.Printf("run_experiment() project_id=%d err=%s\n", project_id, err.Error())
  }

  // send the other fake updates to the project timeline
  go func() {
    // wait to update the timeline status
    time.Sleep(90 * time.Minute)
    
    // send the transform update to arcturus.io project timeline
    proc := exec.Command("curl", 
                         "--insecure", 
                         "-X", "POST", fmt.Sprintf("https://www.arcturus.io/api/projects/%s/activities?access_token=55d28fc5783172b90fea425a2312b95a&key=2", project_id))
    _, err := proc.CombinedOutput()
    if err != nil {
      fmt.Printf("run_experiment() project_id=%d err=%s\n", project_id, err.Error())
    }

    // wait to update the timeline status
    time.Sleep(10 * time.Minute)

    // send the plating update to arcturus.io project timeline
    proc = exec.Command("curl", 
                         "--insecure", 
                         "-X", "POST", fmt.Sprintf("https://www.arcturus.io/api/projects/%s/activities?access_token=55d28fc5783172b90fea425a2312b95a&key=3", project_id))
    _, err = proc.CombinedOutput()
    if err != nil {
      fmt.Printf("run_experiment() project_id=%d err=%s\n", project_id, err.Error())
    }
    
    // wait to update the timeline status
    time.Sleep(5 * time.Minute)
    
    // send the incubating update to arcturus.io project timeline
    proc = exec.Command("curl", 
                         "--insecure", 
                         "-X", "POST", fmt.Sprintf("https://www.arcturus.io/api/projects/%s/activities?access_token=55d28fc5783172b90fea425a2312b95a&key=4", project_id))
    _, err = proc.CombinedOutput()
    if err != nil {
      fmt.Printf("run_experiment() project_id=%d err=%s\n", project_id, err.Error())
    }

    // set the running state
    running = false 
  }()
  
  return err  
}

func run_experiment(project_id string, slot string, genetic_parts string) (error) {
    // set the running state
    running = true

    // turn on streaming
    turn_on_streaming()
    
    // kill the streaming after 5 minutes
    go func() {
      time.Sleep(300 * time.Second)
      turn_off_streaming()
      
      // get the latest video processed by cine.io. It takes some time to encode the video
      // by now we are not getting the final video but adding a previous recorded
      // update the project with the final video
      proc := exec.Command("curl",
                            "--insecure",
                            "-X", "PUT", fmt.Sprintf("https://www.arcturus.io/api/projects/%s?access_token=55d28fc5783172b90fea425a2312b95a&recording_file_name=XJRl3Bsq.20150411T022627.mp4", project_id))
      _, err := proc.CombinedOutput()
      if err != nil {
        fmt.Printf("run_experiment() project_id=%d err=%s\n", project_id, err.Error())
      }                                                     
    }()
    
    // send the assembly update to arcturus.io project timeline
    proc := exec.Command("curl", 
                         "--insecure", 
                         "-X", "POST", fmt.Sprintf("https://www.arcturus.io/api/projects/%s/activities?access_token=55d28fc5783172b90fea425a2312b95a&key=1", project_id))
    _, err := proc.CombinedOutput()
    if err != nil {
      fmt.Printf("run_experiment() project_id=%d err=%s\n", project_id, err.Error())
    }
    
    // run the assembly process. it calls the transforming, plating and incubating
    go func() {
      buf := ""
      fmt.Sprintf(buf, "cmd: python " + filepath.Dir(dir) + "/labcontrol/labcontrol.py -S %s -v -w " + filepath.Dir(dir) + "/labcontrol -A %s -P %s -s assembly_protocol.py", slot, genetic_parts, project_id)
      log.Printf(buf)
      proc = exec.Command("python", filepath.Dir(dir) + "/labcontrol/labcontrol.py", "-S", slot, "-v", "-w", filepath.Dir(dir) + "/labcontrol", "-A", genetic_parts, "-P", project_id, "-s", "assembly_protocol.py")
      proc.Run()
      running = false
    }()
    
    return err        
}

func camera_picture_petri_dish(light_type string) {
  turn_off_streaming()
  
  // remove files
  os.Remove(filepath.Dir(dir) + "/horus-v2/bin/capture.png")
  
  // calling script to take picture
  proc := exec.Command("bash", filepath.Dir(dir) + "/horus-v2/bin/camera-picture-igor.sh", light_type)
  proc.Run()
}

func camera_picture(project_id int, slot string, uv_light bool, light bool) (error) {
  // we cannot use two cameras at the same time
  turn_off_streaming()
  
  // remove files
  os.Remove(filepath.Dir(dir) + "/horus-v2/bin/capture.png")
  
  // run scripts to open the oven, positioning on the grid, open the petri dish
  running = true
  proc := exec.Command("python", filepath.Dir(dir) + "/labcontrol/labcontrol.py", "-S", slot, "-v", "-w", filepath.Dir(dir) + "/labcontrol", "-s", "openOven_openPetriDish_putCamera.py")
  proc.Run()
    
  // take picture with uv_light
  if (uv_light) {
    turn_on_uv_light()
  } 
  
  // take picture with light or no light
  if (light == false) {
    turn_off_light()
  }

  // calling script to take picture
  proc = exec.Command("bash", filepath.Dir(dir) + "/horus-v2/bin/camera-picture.sh")
  proc.Run()
  
  go func() {
    // post picture with curl instead of github.com/ddliu/go-httpclient because I am facing problems with the cacerts from the bbb
    proc = exec.Command("curl", 
                         "--insecure", 
                         "-X", "POST", fmt.Sprintf("https://www.arcturus.io/api/projects/%d/activities?access_token=55d28fc5783172b90fea425a2312b95a&key=5", project_id), 
                         "-F", "content=@"+ filepath.Dir(dir) + "/horus-v2/bin/capture.png")
    _, err := proc.CombinedOutput()

    if err != nil {
      fmt.Printf("camera_picture() project_id=%d err=%s\n", project_id, err.Error())
    }
  }()
  
  // switch uv_light
  if (uv_light) {
    turn_off_uv_light()
  } 
  
  // switch light
  if (light == false) {
    turn_on_light()
  }
  
  // close the  petri dish, turn off the UV light, close the oven, go home
  proc = exec.Command("python", filepath.Dir(dir) + "/labcontrol/labcontrol.py", "-S", slot, "-v", "-w", filepath.Dir(dir) + "/labcontrol", "-s", "closePetriDish_closeOven_goHome.py")
  proc.Run()

  running = false
    
  return nil
}

func turn_on_streaming() (int, error) {
  return exe_cmd(filepath.Dir(dir) + "/horus-v2/bin/camera-streaming.sh")
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

func turn_off_uv_light() (string, error) {
  return serial_cmd("0")
}

func turn_on_uv_light() (string, error) {
  return serial_cmd("1")
}

func turn_off_white_light() (string, error) {
  return serial_cmd("C")
}

func turn_on_white_light() (string, error) {
  return serial_cmd("B")
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

func turn_on_centrifuge() (string, error) {
  return serial_cmd("5")
}

func turn_off_centrifuge() (string, error) {
  return serial_cmd("6")
}

func turn_on_shaker() (string, error) {
  return serial_cmd("7")
}

func turn_off_shaker() (string, error) {
  return serial_cmd("8")
}

func turn_on_gel() (string, error) {
  return serial_cmd("9")
}

func turn_off_gel() (string, error) {
  return serial_cmd("A")
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

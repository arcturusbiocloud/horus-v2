package main

import (
  "github.com/go-martini/martini"
  "github.com/codegangsta/martini-contrib/render"
)

func main() {
  m := martini.Classic()
  m.Use(render.Renderer())
    
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
    
  m.Run()
}

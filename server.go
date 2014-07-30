package main

import (
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/auth"
  "os"
  "log"
  "path/filepath"
  "github.com/jalkoby/martini"
)

func main() {
  SetupUploaders()
  app := martini.New()

  app.Use(martini.Recovery())
  app.Use(martini.Logger())
  app.Use(martini.Static("assets", martini.StaticOptions{Prefix: "assets", }))
  app.Use(martini.Static("uploads", martini.StaticOptions{Prefix: "files"}))

  r := martini.NewRouter()

  r.Group("/files/:id", func(r martini.Router) {
    ctrl := new(FilesCtrl)
    r.Post("", ctrl.Create)
    r.Get("/:dimention", ctrl.Version)
  })

  password := os.Getenv("UPI_PASSWORD")
  if len(password) < 6 {
    panic("Admin password($UPI_PASSWORD) should be at least 6 characters long")
  }
  r.Group("", func(r martini.Router) {
    ctrl := new(ProjectsCtrl)
    r.Get("/", ctrl.Index)

    r.Post("/projects", ctrl.Create)
    r.Delete("/projects/:id", ctrl.Destroy)
  }, render.Renderer(render.Options{Directory: filepath.Join(martini.Root, "templates"), Layout: "layout"}), auth.Basic("admin", password))

  app.Action(r.Handle)
  app.Run()
}

func debugLog(key interface{}, item interface{}) {
  log.Printf("%v: %v\n", key, item)
}

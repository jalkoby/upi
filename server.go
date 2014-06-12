package main

import (
  "github.com/jalkoby/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/auth"
  "os"
  "log"
  "path/filepath"
)

func main() {
  SetupUploaders()
  app := martini.New()

  app.Use(martini.Recovery())
  app.Use(martini.Logger())
  app.Use(martini.Static("assets", martini.StaticOptions{Prefix: "assets", }))
  app.Use(martini.Static("uploads", martini.StaticOptions{Prefix: "files"}))

  r := martini.NewRouter()

  r.Post("/files/:projectId", new(FilesCtrl).Create)

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

func debugLog(str string) {
  log.Println(str)
}

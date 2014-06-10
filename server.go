package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
)

func main() {
  app := martini.New()

  app.Use(martini.Recovery())
  app.Use(martini.Logger())
  app.Use(martini.Static("assets", martini.StaticOptions{Prefix: "assets"}))

  r := martini.NewRouter()

  r.Group("/files", func(r martini.Router) {
    ctrl := new(FilesCtrl)
    r.Get("/:projectId/:id", ctrl.Show)
    r.Post("/:projectId", ctrl.Create)
  })

  r.Group("", func(r martini.Router) {
    ctrl := new(ProjectsCtrl)
    r.Get("/", ctrl.Index)

    r.Post("/projects", ctrl.Create)
    r.Delete("/projects/:id", ctrl.Destroy)
  }, render.Renderer(render.Options{Layout: "layout"}))

  app.Action(r.Handle)
  app.Run()
}

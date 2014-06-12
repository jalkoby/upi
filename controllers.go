package main

import (
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "time"
)

var assetToken = time.Now().Unix()

type ProjectsCtrl struct {}

func (_ ProjectsCtrl) Index(r render.Render) {
  projects, err := AllProjects()
  if err == nil {
    bindings := map[string]interface{}{"projects": projects, "assetToken": assetToken}
    r.HTML(200, "index", bindings)
  } else {
    r.HTML(500, "error", nil)
  }
}

func (_ ProjectsCtrl) Create(request *http.Request, res http.ResponseWriter) (int, string) {
  name := request.PostFormValue("name")
  storage := request.PostFormValue("storage")
  CreateProject(name, storage)
  res.Header().Set("Location", "/")
  return 302, ""
}

func (_ ProjectsCtrl) Destroy(params martini.Params, res http.ResponseWriter) (int, string) {
  err := DestroyProject(params["id"])
  if err == nil {
    res.Header().Set("Location", "/")
    return 204, ""
  } else {
    return 500, ""
  }
}

type FilesCtrl struct {}

func (_ FilesCtrl) Create(params martini.Params, r *http.Request) (int, string) {
  project, err := FindProject(params["projectId"])
  if err != nil { return 404, "not found project" }

  err = r.ParseMultipartForm(1000000)
  if err != nil { return 422, "problem with file parsing" }

  fhs := r.MultipartForm.File["file"]
  if len(fhs) == 0 { return 422, "there is not attached file" }

  file, err := fhs[0].Open()
  defer file.Close()
  if err != nil { return 422, "could not open file" }

  url, err := Upload(file, project)
  if err != nil { return 500, "problem with file storing" }

  return 201, url
}

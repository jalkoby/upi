package main

import (
  "encoding/json"
  "github.com/jalkoby/martini"
  "github.com/martini-contrib/render"
  "net/http"
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

func (_ FilesCtrl) Create(params martini.Params, r *http.Request, w http.ResponseWriter) (int, string) {
  project, err := FindProject(params["id"])
  if err != nil { return 404, "not found project" }

  err = r.ParseMultipartForm(1000000)
  if err != nil { return 422, "problem with file parsing" }

  debugLog("form values", r.MultipartForm.Value)

  urls := map[string]string{}
  for paramName, fs := range r.MultipartForm.File {
    file, err := fs[0].Open()
    if err != nil { return 422, "could not open file" }
    defer file.Close()

    url, err := Upload(file, project)
    if err != nil { return 500, "problem with file storing" }
    urls[paramName] = url
  }

  jsonUrls, _ := json.Marshal(urls)
  w.Header().Set("Content-Type", "application/json")
  return 201, string(jsonUrls)
}

func (_ FilesCtrl) Version(params martini.Params, w http.ResponseWriter) (int, string) {
  file, err := FindFileWithProject(params["id"])
  if err != nil { return 404, "not found file" }
  w.Write(Version(file, params["dimention"]))
  return 200, string(file.Id)
}

package main

import (
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
)

type ProjectsCtrl struct {}

func (_ ProjectsCtrl) Index(r render.Render) {
  projects, err := AllProjects()
  if err == nil {
    r.HTML(200, "index", projects)
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


func (_ FilesCtrl) Show(params martini.Params) (int, string) {
  return 200, params["id"]
}

func (_ FilesCtrl) Create(params martini.Params, r *http.Request) (int, string) {
  project, err := FindProject(params["projectId"])
  if err != nil { return 404, "not found project" }

  err = r.ParseMultipartForm(1000000)
  if err != nil { return 422, "" }

  fhs := r.MultipartForm.File["file"]
  if len(fhs) == 0 { return 422, "" }

  file, err := fhs[0].Open()
  defer file.Close()
  if err != nil { return 422, "" }

  projectToken := project.Token()
  id, err := project.GetFileName()
  if err != nil { return 500, "" }

  if project.IsLocal() {
    err = LocalUpload(file, projectToken, id)
  } else {
    err = S3Upload(file, projectToken, id)
  }

  if err != nil { return 500, "" }
  return 200, "/files/" + projectToken + "/" + id
}

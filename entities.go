package main

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "errors"
)

var Storages = [2]string{"local", "s3"}

type Project struct {
  Id       bson.ObjectId `bson:"_id,omitempty"`
  Name     string        `bson:"name"`
  Storage  string        `bson:"storage"`
}

type File struct {
  Id bson.ObjectId `bson:"_id,omitempty"`
  ProjectId bson.ObjectId `bson:"project_id,omitempty"`
  Project Project
}

func init() {
  c, err := getCollection("files")
  if err != nil { panic(err) }
  defer c.Database.Session.Close()

  index := mgo.Index{
    Key: []string{"project_id"},
    Unique: false,
    DropDups: false,
    Background: true,
    Sparse: true,
  }
  err = c.EnsureIndex(index)
  if err != nil { panic(err) }
}

func AllProjects() ([]Project, error) {
  var projects []Project
  c, err := getCollection("projects")
  if err != nil { return projects, err }
  defer c.Database.Session.Close()

  c.Find(nil).All(&projects)
  return projects, err
}

func CreateProject(name string, storage string) (err error) {
  storageIncluded := func() bool {
    length := len(Storages)
    for i := 0; i < length; i++ { if Storages[i] == storage { return true } }
    return false
  }
  if len(name) == 0 || len(storage) == 0 || !storageIncluded() { return errors.New("invalid entity") }
  c, err :=  getCollection("projects")
  if err != nil { return err }
  defer c.Database.Session.Close()
  err = c.Insert(Project{Name: name, Storage: storage})
  return err
}

func FindProject(_id interface{}) (project Project, err error) {
  c, err := getCollection("projects")
  if err != nil { return project, err }
  defer c.Database.Session.Close()
  var id bson.ObjectId

  switch _id.(type) {
    case bson.ObjectId:
      id = _id.(bson.ObjectId)
    case string:
      id = bson.ObjectIdHex(_id.(string))
    default:
      return project, errors.New("invalid id type")
  }
  err = c.FindId(id).One(&project)
  return project, err
}

func FindFileWithProject(id string) (file File, err error) {
  c, err := getCollection("files")
  if err != nil { return file, err }
  defer c.Database.Session.Close()
  err = c.FindId(bson.ObjectIdHex(id)).One(&file)
  if err != nil { return file, err }
  project, err := FindProject(file.ProjectId)
  file.Project = project
  return file, err
}

func DestroyProject(id string) (err error) {
  c, err := getCollection("projects")
  if err != nil { return err }
  defer c.Database.Session.Close()
  return c.RemoveId(bson.ObjectIdHex(id))
}

func (p Project) IsLocal() bool {
  return p.Storage == Storages[0]
}

func (p Project) GetFileName() (name string, err error) {
  c, err := getCollection("files")
  if err != nil { return name, err }
  defer c.Database.Session.Close()

  f := File{Id: bson.NewObjectId(), ProjectId: p.Id}

  err = c.Insert(f)
  return f.Id.Hex(), err
}

func (p Project) FilesCount() int {
  c, _ := getCollection("files")
  defer c.Database.Session.Close()
  count, _ := c.Find(bson.D{{"project_id", p.Id}}).Count()
  return count
}

func (p Project) Token() string {
  return p.Id.Hex()
}

func (f File) Token() string {
  return f.Id.Hex()
}

func getCollection(name string) (c *mgo.Collection, err error) {
  session, err := mgo.Dial("localhost")
  if err != nil { return c, err }

  return session.DB("upi").C(name), err
}

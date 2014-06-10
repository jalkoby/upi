package main

import (
  "mime/multipart"
  "io"
  "os"
  "strings"
  "bytes"
  "github.com/rlmcpherson/s3gof3r"
)

var rootFolder string
var keys s3gof3r.Keys

func init() {
  rootFolder = os.Getenv("UPI_UPLOAD")
  if len(rootFolder) == 0 { panic("Please specify $UPI_UPLOAD variable") }

  if !strings.HasSuffix(rootFolder, "/") { rootFolder += "/" }

  keys, err := s3gof3r.EnvKeys()
  if err != nil { panic(err) }
  _ = keys
}

func LocalUpload(file multipart.File, projectToken string, id string) (err error) {
  var filePath = bytes.NewBufferString(rootFolder)
  filePath.WriteString(projectToken)
  filePath.WriteString("/")

  _, err = os.Stat(filePath.String())
  if os.IsNotExist(err) { os.Mkdir(filePath.String(), os.FileMode(0775)) }

  filePath.WriteString(id)
  dst, err := os.Create(filePath.String())
  defer dst.Close()
  if err != nil { return err }

  _, err = io.Copy(dst, file)
  return err
}

func S3Upload(file multipart.File, projectToken string, id string) (err error) {
  bucket := s3gof3r.New("", keys).Bucket(projectToken)
  s3Writer, err := bucket.PutWriter(id, nil, nil)
  if err != nil { return err }
  defer s3Writer.Close()

  _, err = io.Copy(s3Writer, file);
  return err
}

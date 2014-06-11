package main

import (
  "github.com/jalkoby/s3gof3r"
  "io"
  "fmt"
  "net/http"
  "mime/multipart"
  "os"
)

var configs = map[string]string{
  "localFolder": os.Getenv("UPI_UPLOAD"),
  "localHost": os.Getenv("UPI_PUBLIC_HOST"),
  "awsBucket": os.Getenv("AWS_BUCKET"),
}

func SetupUploaders() {
  bucket := getBucket()
  bucket.Create("", false, nil)

  for key, value := range configs {
    if len(value) == 0 {
      panic(fmt.Sprintf("%v is not configured", key))
    }
  }
}

func Upload(file multipart.File, project Project) (string, error) {
  token := project.Token()
  id, err := project.GetFileName()
  if err != nil { return "", err }

  if project.IsLocal() {
    return localUpload(file, token, id)
  } else {
    return s3Upload(file, token, id)
  }
}

func localUpload(file multipart.File, token string, id string) (string, error) {
  var path = fmt.Sprintf("%v/%v", configs["localFolder"], token)

  _, err := os.Stat(path)
  if os.IsNotExist(err) { os.Mkdir(path, os.FileMode(0775)) }

  path = fmt.Sprintf("%v/%v", path, id)
  dst, err := os.Create(path)
  defer dst.Close()
  if err != nil { return "", err }

  _, err = io.Copy(dst, file)
  if err != nil { return "", err }

  url := fmt.Sprintf("%v/%v/%v", configs["localHost"], token, id)
  return url, nil
}

func s3Upload(file multipart.File, token string, id string) (string, error) {
  bucket := getBucket()
  path := fmt.Sprintf("%v/%v", token, id)
  headers := make(http.Header)
  headers.Add("x-amz-acl", "public-read")
  s3Writer, err := bucket.PutWriter(path, headers, nil)
  if err != nil { return "", err }
  defer s3Writer.Close()

  _, err = io.Copy(s3Writer, file);
  if err != nil { return "", err }

  url := fmt.Sprintf("http://%v/%v/%v/%v", bucket.S3.Domain, bucket.Name, token, id)
  return url, nil
}

func getBucket() s3gof3r.Bucket {
  keys, err := s3gof3r.EnvKeys()
  if err != nil { panic(err) }
  return *s3gof3r.New("", keys).Bucket(configs["awsBucket"])
}

package main

import (
  "bytes"
  "github.com/jalkoby/s3gof3r"
  "github.com/gographics/imagick/imagick"
  "io"
  "io/ioutil"
  "fmt"
  "net/http"
  "mime/multipart"
  "os"
  "strconv"
  "strings"
)

var configs = map[string]string{
  "localFolder": os.Getenv("UPI_UPLOAD"),
  "localHost": os.Getenv("UPI_PUBLIC_HOST"),
  "awsBucket": os.Getenv("AWS_BUCKET"),
}

type Size struct {
  Width uint
  Height uint
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

func Version(file File, dimention string) []byte {
  generator := func() []byte {
    imagick.Initialize()
    defer imagick.Terminate()
    mw := imagick.NewMagickWand()
    defer mw.Destroy()
    if file.Project.IsLocal() {
      mw.ReadImage(localPath(file))
    } else {
      mw.ReadImageBlob(s3FileContent(file))
    }
    original := Size{Width: mw.GetImageWidth(), Height: mw.GetImageHeight()}
    size := parseSize(dimention, original)
    mw.ResizeImage(size.Width, size.Height, imagick.FILTER_LANCZOS, 1)
    return mw.GetImageBlob()
  }
  return versionWithCache(file, dimention, generator)
}

func versionWithCache(file File, dimention string, generator func() []byte) []byte {
  cacheDir := fmt.Sprintf("%v/cache/%v", configs["localFolder"], file.Token())
  os.MkdirAll(cacheDir, 0775)
  path := fmt.Sprintf("%v/%v", cacheDir, dimention)
  _, err := os.Stat(path)
  if os.IsNotExist(err) {
    debugLog("dont use case", path)
    content := generator()
    ioutil.WriteFile(path, content, 0664)
    return content
  } else {
    debugLog("use case", path)
    content, _ := ioutil.ReadFile(path)
    return content
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

  url := fmt.Sprintf("%v/files/%v/%v", configs["localHost"], token, id)
  return url, nil
}

func localPath(file File) string {
  return fmt.Sprintf("%v/%v/%v", configs["localFolder"], file.Project.Token(), file.Token())
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

func s3FileContent(file File) []byte {
  bucket := getBucket()
  path := fmt.Sprintf("%v/%v", file.Project.Token(), file.Token())
  s3Reader, _, err := bucket.GetReader(path, nil)
  if err != nil { panic(err) }
  buf := new(bytes.Buffer)
  buf.ReadFrom(s3Reader)
  return buf.Bytes()
}

func getBucket() s3gof3r.Bucket {
  keys, err := s3gof3r.EnvKeys()
  if err != nil { panic(err) }
  return *s3gof3r.New("", keys).Bucket(configs["awsBucket"])
}

func parseSize(dimention string, original Size) Size {
  parts := strings.Split(dimention, "x")
  var width, height uint
  if len(parts) > 0 {
    w, _ := strconv.ParseInt(parts[0], 0, 0)
    width = uint(w)
  }
  if len(parts) > 1 {
    h, _ := strconv.ParseInt(parts[1], 0, 0)
    height = uint(h)
  }

  if (width == 0) && (height == 0) {
    width = original.Width
    height = original.Height
  } else if width == 0 {
    width = uint(float32(original.Width * height)/float32(original.Height))
  } else if height == 0 {
    height = uint(float32(original.Height * width)/float32(original.Width))
  }

  return Size{Width: width, Height: height }
}

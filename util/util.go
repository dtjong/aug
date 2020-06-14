package util

import (
  "os"
  "strings"
  "path/filepath"
)

func FilterExt(arr []os.FileInfo, ext string) []os.FileInfo {
  filtered := make([]os.FileInfo, 0)

  for _, item := range arr {
    filename := item.Name()
    fileExt := filepath.Ext(filename)
    if strings.ToLower(fileExt) == strings.ToLower(ext) {
      filtered = append(filtered, item)
    }
  }

  return filtered
}

func Check(e error) {
  if e != nil {
    panic(e)
  }
}

func TrimExt(filename string) string {
  extension := filepath.Ext(filename)
  return filename[0: len(filename) - len(extension)]
}

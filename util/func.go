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

package util

import (
  "os"
  "strings"
  "path/filepath"
  "math/rand"
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

func RandInt(maxValue int) int {
  return rand.Int() % maxValue
}

func RandIntExp(lambda float64) int {
  return int(rand.ExpFloat64() * lambda)
}

func RandFloatExp(lambda float64) float64 {
  return rand.ExpFloat64() * lambda
}

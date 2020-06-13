package main

import (
  "fmt"

  "aug/imagedata"
  "aug/readdata"
)

const loadDir = "Data/test/"
const writeDir = "Data/dest/"

const labelFormat = imagedata.XML

func main() {
  images := readdata.ReadData(loadDir, labelFormat)
  fmt.Println("Done augmenting!")
}

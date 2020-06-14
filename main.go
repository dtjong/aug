package main

import (
  "fmt"

  "aug/imagedata"
  "aug/readdata"
  "aug/writedata"
  "aug/augment"
)

const loadDir = "Data/test/"
const writeDir = "Data/dest/"

const labelFormat = imagedata.XML
const outputFormat = imagedata.YOLO

const augmentationsPerImage = 10

func main() {
  images := readdata.ReadData(loadDir, labelFormat)
  augmented := augment.Augment(images, augmentationsPerImage)
  writedata.WriteData(writeDir, augmented, outputFormat)
  fmt.Println("Done augmenting!")
}

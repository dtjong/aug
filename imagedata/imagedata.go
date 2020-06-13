package imagedata

import (
  "gocv.io/x/gocv"
)

type ImageData struct {
  BBoxes []BBox
  Image gocv.Mat
  Name string
}

type BBox struct {
  Class int
  XMin float64
  YMin float64
  XMax float64
  YMax float64
}

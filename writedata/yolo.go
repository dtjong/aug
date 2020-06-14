package writedata

import (
  "os"
  "fmt"

  "gocv.io/x/gocv"

  "aug/imagedata"
  "aug/util"
)

const imagesDir = "images/"
const labelsDir = "labels/"
const lists = "lists.txt"
const imageExt = ".jpg"
const labelExt = ".txt"

type yoloBBox struct {
  Class int
  X float64
  Y float64
  Width float64
  Height float64
}

func WriteDataYOLO(writeDir string, images []imagedata.ImageData) {
  err := os.Mkdir(writeDir + imagesDir, 0755)
  util.Check(err)

  err = os.Mkdir(writeDir + labelsDir, 0755)
  util.Check(err)

  cwd, err := os.Getwd()
  util.Check(err)

  // Writing lists.txt
  lists, err := os.Create(writeDir + lists)
  util.Check(err)
  defer lists.Close()

  for _, image := range images {
    imagePath := cwd + "/" + writeDir + imagesDir + image.Name + imageExt
    lists.WriteString(imagePath + "\n")

    // Write image
    gocv.IMWrite(imagePath, image.Image)

    // Write labels
    width := image.Image.Cols()
    height := image.Image.Rows()

    label, err := os.Create(writeDir + labelsDir + image.Name + labelExt)
    util.Check(err)

    for _, box := range image.BBoxes {
      yoloBBox := convertYoloBBox(width, height, box)
      label.WriteString(yoloBBox.fmt())
    }
    label.Close()
  }
}

func convertYoloBBox(width, height int, bbox imagedata.BBox) yoloBBox {
  boxWidth := (bbox.XMax - bbox.XMin) / float64(width)
  boxHeight := (bbox.YMax - bbox.YMin) / float64(height)
  x := boxWidth / 2 + bbox.XMin / float64(width)
  y := boxHeight / 2 + bbox.YMin / float64(height)
  
  return yoloBBox {
    bbox.Class,
    x,
    y,
    boxWidth,
    boxHeight,
  }
}

func (bbox yoloBBox) fmt() string {
  return fmt.Sprintf("%d %f %f %f %f\n",
    bbox.Class, bbox.X, bbox.Y, bbox.Width, bbox.Height)
}

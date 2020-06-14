package readdata

import (
  "encoding/xml"
  "os"
  "io/ioutil"

  "gocv.io/x/gocv"

  "aug/imagedata"
  "aug/util"
)

type annotation struct {
  XMLName xml.Name `xml:"annotation"`
  Boundingboxes []bbox `xml:"object"`
}

type bbox struct {
  XMLName xml.Name `xml:"object"`
  XMin int `xml:"bndbox>xmin"`
  YMin int `xml:"bndbox>ymin"`
  XMax int `xml:"bndbox>xmax"`
  YMax int `xml:"bndbox>ymax"`
}

func ReadDataXml(dirname string, files []os.FileInfo) []imagedata.ImageData {
  labels := util.FilterExt(files, ".xml")
  images := util.FilterExt(files, ".jpg")

  data := make([]imagedata.ImageData, 0)

  for i := 0; i < len(labels); i++ {
    imdata := imagedata.ImageData {
      Image: gocv.IMRead(dirname + images[i].Name(), gocv.IMReadUnchanged),
      BBoxes: readXml(dirname + labels[i].Name()),
      Name: util.TrimExt(images[i].Name()),
    }

    data = append(data, imdata)
  }

  return data
}

func readXml(filepath string) []imagedata.BBox {
  bboxes := make([]imagedata.BBox, 0)

  // Read xml into struct
  file, err := os.Open(filepath)

  if err != nil {
    panic(err)
  }

  defer file.Close()

  data, _ := ioutil.ReadAll(file)

  var annot annotation
  xml.Unmarshal(data, &annot)

  // xml struct to bbox
  for _, box := range annot.Boundingboxes {
    bboxes = append(bboxes, imagedata.BBox {
      Class: 0,
      XMin: float64(box.XMin),
      YMin: float64(box.YMin),
      XMax: float64(box.XMax),
      YMax: float64(box.YMax),
    })
  }

  return bboxes
}

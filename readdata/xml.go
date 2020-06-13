package readdata

import (
  "fmt"
  "encoding/xml"
  "os"
  "io/ioutil"

  "aug/imagedata"
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

func ReadXml(filepath string) []imagedata.BBox {
  bboxes := make([]imagedata.BBox, 0)

  // Read xml into struct
  file, err := os.Open(filepath)

  if err != nil {
    fmt.Println("Error opening file " + filepath)
    return bboxes
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
      YMax: float64(box.XMax),
    })
  }

  return bboxes
}

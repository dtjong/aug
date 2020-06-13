package readdata

import (
  "io/ioutil"
  "os"

  "gocv.io/x/gocv"

  "aug/imagedata"
  "aug/util"
)

func ReadData(dirname string, labelFormat imagedata.LabelFormat) []imagedata.ImageData {
  files, err := ioutil.ReadDir(dirname)

  if err != nil {
    os.Exit(0)
  }

  labels := util.FilterExt(files, ".xml")
  images := util.FilterExt(files, ".jpg")

  data := make([]imagedata.ImageData, 0)

  for i := 0; i < len(labels); i++ {
    imdata := imagedata.ImageData {
      Image: gocv.IMRead(dirname + images[i].Name(), gocv.IMReadUnchanged),
      BBoxes: ReadXml(dirname + labels[i].Name()),
      Name: images[i].Name(),
    }

    data = append(data, imdata)
  }

  return data
}

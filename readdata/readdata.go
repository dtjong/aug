package readdata

import (
  "io/ioutil"

  "aug/imagedata"
)

func ReadData(dirname string, labelFormat imagedata.LabelFormat) []imagedata.ImageData {
  files, err := ioutil.ReadDir(dirname)

  if err != nil {
    panic(err)
  }
  
  switch labelFormat {
    case imagedata.XML:
      return ReadDataXml(dirname, files)
    default:
      return make([]imagedata.ImageData, 0)
  }
}


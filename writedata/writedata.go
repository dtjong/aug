package writedata

import (
  "os"

  "aug/imagedata"
  "aug/util"
)

func WriteData(
    writeDir string, 
    images []imagedata.ImageData, 
    labelFormat imagedata.LabelFormat) {
  
  err := os.RemoveAll(writeDir)
  util.Check(err)
  err = os.Mkdir(writeDir, 0755)
  util.Check(err)

  switch labelFormat {
    case imagedata.YOLO:
      WriteDataYOLO(writeDir, images)
  }
}


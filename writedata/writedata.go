package writedata

import (
  "os"
  "fmt"

  "aug/imagedata"
  "aug/util"
)

func WriteData(
    writeDir string, 
    images []imagedata.ImageData, 
    labelFormat imagedata.LabelFormat) {
  fmt.Println("Writing files...")
  
  err := os.RemoveAll(writeDir)
  util.Check(err)
  err = os.Mkdir(writeDir, 0755)
  util.Check(err)

  switch labelFormat {
    case imagedata.YOLO:
      WriteDataYOLO(writeDir, images)
  }
}


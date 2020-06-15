package augment

import (
  "image"
  "strconv"
  "math/rand"

  "gocv.io/x/gocv"

  "aug/imagedata"
  "aug/util"
)

var Augmentations = [...]func(imagedata.ImageData) imagedata.ImageData {
  blur,
  vflip,
  hflip,
  noise,
}

/*
Augmentations take in the image data and change the it and the bounding 
box if necessary, along with the name.

Augmentations must include a degree of randomness in it, and include that
in the outputted name.
*/

func blur(img imagedata.ImageData) imagedata.ImageData {
  if rand.Float32() < .4 {
    return img
  }

  level := util.RandIntExp(1.2) * 2 + 1

  gocv.GaussianBlur(img.Image, &(img.Image), image.Pt(level, level), 
    0, 0, gocv.BorderDefault)
  img.Name = img.Name + "b" + strconv.Itoa(level)

  return img
}

func vflip(img imagedata.ImageData) imagedata.ImageData {
  if rand.Float32() < .5 {
    return img
  }

  gocv.Flip(img.Image, &(img.Image), 1)
  img.Name = img.Name + "h"

  width := img.Image.Cols()

  // Flip bounding boxes
  imagedata.Apply(img.BBoxes, func(bbox imagedata.BBox) imagedata.BBox {
    xMin := float64(width) - bbox.XMax
    xMax := float64(width) - bbox.XMin
    bbox.XMin = xMin
    bbox.XMax = xMax
    return bbox
  })

  return img
}

func hflip(img imagedata.ImageData) imagedata.ImageData {
  if rand.Float32() < .75 {
    return img
  }

  gocv.Flip(img.Image, &(img.Image), 0)
  img.Name = img.Name + "v"

  height := img.Image.Rows()

  // Flip bounding boxes
  imagedata.Apply(img.BBoxes, func(bbox imagedata.BBox) imagedata.BBox {
    yMin := float64(height) - bbox.YMax
    yMax := float64(height) - bbox.YMin
    bbox.YMin = yMin
    bbox.YMax = yMax
    return bbox
  })

  return img
}

func noise(img imagedata.ImageData) imagedata.ImageData {
  level := rand.Float32() / 10
  
  height := img.Image.Rows()
  width := img.Image.Cols()

  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      pix := imagedata.GetVecbAt(img.Image, x, y)
      if rnd := rand.Float32(); rnd < level {
        pix.Add(int(rand.Float32() * 500 - 250), len(pix))
      }
      pix.SetVecbAt(img.Image, x, y)
    }
  }

  img.Name = img.Name + "n" + strconv.Itoa(int(level * 1000))

  return img
}

package augment

import (
  "aug/imagedata"
)

/*
Augment takes a slice of imagedata and appends augmentations to the slice.

Augmentations will have a new name descriptive of the augmentation levels, 
and contain bounding boxes and an image matrix to reflect it.
*/
func Augment(images []imagedata.ImageData, augmentationsPerImage int) []imagedata.ImageData {
  augmentedImages := make([]imagedata.ImageData, 0)

  // For each image, create n augmentations, and apply each augmentation
  for _, image := range images {
    for i := 0; i < augmentationsPerImage; i++ {
      augmentedImage := image.Clone()
      for _, augmentation := range Augmentations {
        augmentedImage = augmentation(augmentedImage)
      }
    }
  }

  return append(images, augmentedImages...)
}

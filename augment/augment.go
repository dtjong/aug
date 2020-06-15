package augment

import (
  "fmt"
  "sync"

  "aug/imagedata"
)

/*
Augment takes a slice of imagedata and appends augmentations to the slice.

Augmentations will have a new name descriptive of the augmentation levels, 
and contain bounding boxes and an image matrix to reflect it.
*/
func Augment(images []imagedata.ImageData, augmentationsPerImage int) []imagedata.ImageData {
  augmentedImages := make([]imagedata.ImageData, 0)

  mutex := &sync.Mutex{}
  wg := sync.WaitGroup{}


  // For each image, create n augmentations, and apply each augmentation
  for _, image := range images {
    fmt.Println("Augmenting " + image.Name)
    for i := 0; i < augmentationsPerImage; i++ {
      wg.Add(1)
      go func() {
        augmentedImage := image.Clone()
        for _, augmentation := range Augmentations {
          augmentedImage = augmentation(augmentedImage)
        }
        mutex.Lock()
        augmentedImages = append(augmentedImages, augmentedImage)
        mutex.Unlock()
        wg.Done()
      }()
    }
  }

  wg.Wait()
  return append(images, augmentedImages...)
}


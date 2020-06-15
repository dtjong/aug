package imagedata

import (
  "math"

  "gocv.io/x/gocv"
)

type ImageData struct {
  BBoxes []BBox
  Image gocv.Mat
  Name string
}

type BBox struct {
  Class int
  XMin float64
  YMin float64
  XMax float64
  YMax float64
}

type Vecb []uint8

func (data ImageData) Clone() ImageData {
  bboxesCopy := make([]BBox, len(data.BBoxes))
  copy(bboxesCopy, data.BBoxes)
  return ImageData {
    BBoxes: bboxesCopy,
    Image: data.Image.Clone(),
    Name: data.Name,
  }
}

func GetVecbAt(m gocv.Mat, row int, col int) Vecb {
	ch := m.Channels()
	v := make(Vecb, ch)

	for c := 0; c < ch; c++ {
		v[c] = m.GetUCharAt(row, col*ch+c)
	}

	return v
}

func (v Vecb) SetVecbAt(m gocv.Mat, row int, col int) {
	ch := m.Channels()

	for c := 0; c < ch; c++ {
		m.SetUCharAt(row, col*ch+c, v[c])
	}
}

// Adds value to all pixel values, and caps at max value
func (pix Vecb) Add(val, ch int) {
	for c := 0; c < ch; c++ {
    pix[c] = uint8(math.Max(0, math.Min(math.MaxUint8, float64(val + int(pix[c])))))
  }
}

func Apply(bboxes []BBox, f func(BBox) BBox) {
  for i, box := range bboxes {
    bboxes[i] = f(box)
  }
}

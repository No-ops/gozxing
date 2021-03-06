package testutil

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

func ExpandBitMatrix(src *gozxing.BitMatrix, factor int) *gozxing.BitMatrix {
	dst, _ := gozxing.NewBitMatrix(src.GetWidth()*factor, src.GetHeight()*factor)
	for j := 0; j < src.GetHeight(); j++ {
		y := j * factor
		for i := 0; i < src.GetWidth(); i++ {
			x := i * factor
			if src.Get(i, j) {
				dst.SetRegion(x, y, factor, factor)
			}
		}
	}
	return dst
}

func NewBitArrayFromString(str string) *gozxing.BitArray {
	arr := gozxing.NewBitArray(len(str))
	for i, c := range str {
		if c == '1' {
			arr.Set(i)
		}
	}
	return arr
}

func NewBinaryBitmapFromBitMatrix(matrix *gozxing.BitMatrix) *gozxing.BinaryBitmap {
	src := newTestBitMatrixSource(matrix)
	binarizer := gozxing.NewHybridBinarizer(src)
	bmp, _ := gozxing.NewBinaryBitmap(binarizer)
	return bmp
}

func NewBinaryBitmapFromFile(filename string) *gozxing.BinaryBitmap {
	file, _ := os.Open(filename)
	img, _, _ := image.Decode(file)
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	return bmp
}

type testBitMatrixSource struct {
	gozxing.LuminanceSourceBase
	matrix *gozxing.BitMatrix
}

func newTestBitMatrixSource(matrix *gozxing.BitMatrix) gozxing.LuminanceSource {
	return &testBitMatrixSource{
		gozxing.LuminanceSourceBase{matrix.GetWidth(), matrix.GetHeight()},
		matrix,
	}
}

func (this *testBitMatrixSource) GetRow(y int, row []byte) ([]byte, error) {
	for x := 0; x < this.matrix.GetWidth(); x++ {
		if this.matrix.Get(x, y) {
			row[x] = 0
		} else {
			row[x] = 255
		}
	}
	return row, nil
}

func (this *testBitMatrixSource) GetMatrix() []byte {
	width := this.GetWidth()
	height := this.GetHeight()
	matrix := make([]byte, width*height)
	for y := 0; y < height; y++ {
		offset := y * width
		for x := 0; x < width; x++ {
			if !this.matrix.Get(x, y) {
				matrix[offset+x] = 255
			}
		}
	}
	return matrix
}

func (this *testBitMatrixSource) Invert() gozxing.LuminanceSource {
	return gozxing.LuminanceSourceInvert(this)
}

func (this *testBitMatrixSource) String() string {
	return gozxing.LuminanceSourceString(this)
}

type DummyGridSampler struct{}

func (s DummyGridSampler) SampleGrid(image *gozxing.BitMatrix, dimensionX, dimensionY int,
	p1ToX, p1ToY, p2ToX, p2ToY, p3ToX, p3ToY, p4ToX, p4ToY float64,
	p1FromX, p1FromY, p2FromX, p2FromY, p3FromX, p3FromY, p4FromX, p4FromY float64) (*gozxing.BitMatrix, error) {
	return nil, errors.New("dummy sampler")
}

func (s DummyGridSampler) SampleGridWithTransform(image *gozxing.BitMatrix,
	dimensionX, dimensionY int, transform *common.PerspectiveTransform) (*gozxing.BitMatrix, error) {
	return nil, errors.New("dummy sampler")
}

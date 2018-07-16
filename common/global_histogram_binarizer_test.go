package common

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

type testLuminanceSource struct {
	gozxing.LuminanceSourceBase
}

func newTestLuminanceSource(size int) *testLuminanceSource {
	return &testLuminanceSource{
		gozxing.LuminanceSourceBase{size, size},
	}
}
func (this *testLuminanceSource) GetRow(y int, row []byte) []byte {
	width := this.GetWidth()
	for i := 0; i < width; i++ {
		if (y+i)%2 == 0 {
			row[i] = 10 + byte(50*i/width)
		} else {
			row[i] = 250 - byte(50*i/width)
		}
	}
	return row
}
func (this *testLuminanceSource) GetMatrix() []byte {
	width := this.GetWidth()
	height := this.GetHeight()
	matrix := make([]byte, width*height)
	for y := 0; y < height; y++ {
		this.GetRow(y, matrix[width*y:])
	}
	return matrix
}
func (this *testLuminanceSource) Invert() gozxing.LuminanceSource {
	return gozxing.LuminanceSourceInvert(this)
}
func (this *testLuminanceSource) String() string {
	return gozxing.LuminanceSourceString(this)
}

type testBlackSource struct {
	gozxing.LuminanceSourceBase
}

func newTestBlackSource(size int) *testBlackSource {
	return &testBlackSource{
		gozxing.LuminanceSourceBase{size, size},
	}
}
func (this *testBlackSource) GetRow(y int, row []byte) []byte {
	for i := 0; i < this.GetWidth(); i++ {
		row[i] = 0
	}
	return row
}
func (this *testBlackSource) GetMatrix() []byte {
	size := this.GetWidth() * this.GetHeight()
	matrix := make([]byte, size)
	return matrix
}
func (this *testBlackSource) Invert() gozxing.LuminanceSource {
	return gozxing.LuminanceSourceInvert(this)
}
func (this *testBlackSource) String() string {
	return gozxing.LuminanceSourceString(this)
}

func TestGlobalHistgramBinarizer(t *testing.T) {
	size := 32
	src := newTestLuminanceSource(size)
	ghb := NewGlobalHistgramBinarizer(src)

	if s := ghb.GetLuminanceSource(); s != src {
		t.Fatalf("GetLuminanceSource = %p, expect %p", s, src)
	}
	if w, h := ghb.GetWidth(), ghb.GetHeight(); w != size || h != size {
		t.Fatalf("GetWidth,GetHeight = %v,%v, expect %v,%v", w, h, size, size)
	}
}

func TestGlobalHistgramBinarizer_estimateBlackPoint(t *testing.T) {
	g := GlobalHistogramBinarizer{}

	// single peak
	buckets := []int{0, 0, 0, 15, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	_, e := g.estimateBlackPoint(buckets)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("estimateBlackPoint must be NotFoundException, %T", e)
	}

	buckets = []int{0, 0, 0, 15, 12, 12, 5, 14, 16, 19, 20, 18, 0, 0, 0, 0}
	valley := 6 << LUMINANCE_SHIFT
	r, e := g.estimateBlackPoint(buckets)
	if e != nil {
		t.Fatalf("estimateBlackPoint returns error, %v", e)
	}
	if r != valley {
		t.Fatalf("estimateBlackPoint = %v, expect %v", r, valley)
	}
}

func TestGlobalHistgramBinarizer_GetBlackRow(t *testing.T) {
	src := newTestLuminanceSource(16)
	ghb := NewGlobalHistgramBinarizer(src)

	r, e := ghb.GetBlackRow(1, nil)
	if e != nil {
		t.Fatalf("GetBlackRow returns error, %v", e)
	}
	expect := " .X.X.X.X .X.X.X.."
	if r.String() != expect {
		t.Fatalf("GetBlackRow = \"%v\", expect \"%v\"", r, expect)
	}

	// white image
	ghb = ghb.CreateBinarizer(newTestBlackSource(16))
	_, e = ghb.GetBlackRow(0, r)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("GetBlackRow must be NotFoundException, %T", e)
	}

	// small image
	ghb = ghb.CreateBinarizer(newTestLuminanceSource(2))
	r, e = ghb.GetBlackRow(0, nil)
	if e != nil {
		t.Fatalf("GetBlackRow returns error, %v", e)
	}
	expect = " X."
	if r.String() != expect {
		t.Fatalf("GetBlackRow = \"%v\", expect \"%v\"", r, expect)
	}
}

func TestGlobalHistgramBinarizer_GetBlackMatrix(t *testing.T) {
	ghb := NewGlobalHistgramBinarizer(newTestLuminanceSource(0))
	_, e := ghb.GetBlackMatrix()
	if e == nil {
		t.Fatalf("GetBlackMatrix must be error")
	}

	ghb = NewGlobalHistgramBinarizer(newTestBlackSource(16))
	_, e = ghb.GetBlackMatrix()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("GetBlackMatrix must be NotFoundException, %T", e)
	}

	src := newTestLuminanceSource(16)
	rawmatrix := src.GetMatrix()
	ghb = NewGlobalHistgramBinarizer(src)
	m, e := ghb.GetBlackMatrix()
	if e != nil {
		t.Fatalf("GetBlackMatrix returns error, %v", e)
	}
	for w := 0; w < m.GetWidth(); w++ {
		for h := 0; h < m.GetHeight(); h++ {
			expect := rawmatrix[w+m.GetHeight()*h] < 128
			if r := m.Get(w, h); r != expect {
				t.Fatalf("GetBlackMatrix [%v,%v] is %v, expect %v", w, h, r, expect)
			}
		}
	}
}
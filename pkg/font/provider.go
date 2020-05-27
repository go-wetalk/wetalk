package font

import (
	"image"
	"io/ioutil"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

func GetContext() (*freetype.Context, error) {
	b, err := ioutil.ReadFile(FontPath)
	if err != nil {
		return nil, err
	}

	f, err := freetype.ParseFont(b)
	if err != nil {
		return nil, err
	}

	fc := freetype.NewContext()
	fc.SetHinting(font.HintingNone)
	fc.SetFont(f)
	fc.SetSrc(image.Black)

	return fc, nil
}

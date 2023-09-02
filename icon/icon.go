package icon

import (
	"bufio"
	"bytes"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"pc-sensors-tray/types"
)

var (
	dpi      = 160.0
	fontfile = "/usr/share/fonts/truetype/ubuntu/Ubuntu-B.ttf"
	size     = 10.0
)

type IconService struct {
	face  font.Face
	cache map[string][]byte
}

func NewIconService() IconService {
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Fatal(err)
	}
	parsedFont, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(parsedFont, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	cache := make(map[string][]byte)
	return IconService{face, cache}
}

func (svc *IconService) GetIcon(result types.Result) ([]byte, error) {
	cached, ok := svc.cache[result.Icon()]
	if ok {
		return cached, nil
	}

	fgCol := colornames.Red
	fgCol.G = uint8(255.0 - 255.0*result.Intensity())
	fgCol.B = uint8(255.0 - 255.0*result.Intensity())
	fg, bg := image.NewUniform(fgCol), image.NewUniform(color.RGBA{0x12, 0x12, 0x12, 0xff})

	const imgW, imgH = 32, 32
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	d := &font.Drawer{
		Dst:  rgba,
		Src:  fg,
		Face: svc.face,
	}
	y := int(math.Ceil(size * dpi / 72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) - d.MeasureString(result.Icon())) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(result.Icon())

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err := png.Encode(writer, rgba)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	bs := b.Bytes()
	svc.cache[result.Icon()] = bs
	return bs, nil
}

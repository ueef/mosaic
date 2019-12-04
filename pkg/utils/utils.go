package utils

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"runtime"
	"sync"
)

func Rotate270(i image.Image) image.Image {
	src := ConvertToRGBA(i)
	b := src.Bounds()
	dst := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: b.Min.Y,
			Y: b.Min.X,
		},
		Max: image.Point{
			X: b.Max.Y,
			Y: b.Max.X,
		},
	})
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			s := y*src.Stride + x*4
			d := (b.Max.X-x-1)*dst.Stride + y*4
			copy(dst.Pix[d:d+4], src.Pix[s:s+4])
		}
	}

	return dst
}

func Rotate180(i image.Image) image.Image {
	src := ConvertToRGBA(i)
	b := src.Bounds()
	dst := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: b.Min.Y,
			Y: b.Min.X,
		},
		Max: image.Point{
			X: b.Max.Y,
			Y: b.Max.X,
		},
	})
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			s := y*src.Stride + x*4
			d := (b.Max.X-x-1)*dst.Stride + (b.Max.Y-y-1)*4
			copy(dst.Pix[d:d+4], src.Pix[s:s+4])
		}
	}

	return src
}

func Rotate90(i image.Image) image.Image {
	src := ConvertToRGBA(i)
	b := src.Bounds()
	dst := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: b.Min.Y,
			Y: b.Min.X,
		},
		Max: image.Point{
			X: b.Max.Y,
			Y: b.Max.X,
		},
	})
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			s := y*src.Stride + x*4
			d := x*dst.Stride + (b.Max.Y-y-1)*4
			copy(dst.Pix[d:d+4], src.Pix[s:s+4])
		}
	}

	return dst
}

func ConvertToRGBA(src image.Image) (dst *image.RGBA) {
	rgba, ok := src.(*image.RGBA)
	if !ok {
		rgba = image.NewRGBA(src.Bounds())
		draw.Draw(rgba, src.Bounds(), src, src.Bounds().Min, draw.Src)
	}

	return rgba
}

func Line(ls, le int, f func(s, e int)) {
	p := runtime.GOMAXPROCS(0)
	ll := le - ls
	ps := ll / p
	if p <= 1 || ps <= p {
		f(0, ll)
		return
	}

	wg := sync.WaitGroup{}
	pn := ll
	for pn > 0 {
		s := ls + pn - ps
		if s < 0 {
			s = 0
		}

		e := s + pn
		if e > le {
			e = le
		}

		pn -= ps
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(s, e)
		}()
	}
	wg.Wait()
}

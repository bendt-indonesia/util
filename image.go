package util

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"github.com/bendt-indonesia/util/enum"

)

type ImageReader struct {
	Bounds   image.Rectangle
	FileName string
	Format   string
	Image    *image.Image
	Height   int
	Width    int
	FileSize int64
}

type ImageOutput struct {
	Source     *image.Image
	SourceRGBA *image.RGBA
	Fit        *enum.ImageFit
	FilePath   string
	Format     string
	Height     *int
	Width      *int
	Quality    int
}

type ImageResizeOptions struct {
	Height        uint
	Width         uint
	MaintainRatio bool
	Quality       int
}

func (i *ImageOutput) SetDefault() {
	if i.FilePath == "" {
		i.SetFilePath(nil)
	}

	i.Quality = jpeg.DefaultQuality
}

func (i *ImageOutput) SetFilePath(overridePath *string) {
	if overridePath != nil {
		i.FilePath = *overridePath
		i.Format = GetImageExt(i.FilePath)
	} else {
		i.Format = "jpg"
		i.FilePath = RandomTimestampStr() + "." + i.Format
	}
}

func (i *ImageOutput) SaveAs() {
	if i.Source == nil && i.SourceRGBA == nil {
		return
	}

	i.SetDefault()

	newImg, err := os.Create(i.FilePath)
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}

	if i.Format == "png" {
		if i.Source != nil {
			png.Encode(newImg, *i.Source)
		} else if i.SourceRGBA != nil {
			png.Encode(newImg, i.SourceRGBA)
		}
	} else {
		if i.Source != nil {
			jpeg.Encode(newImg, *i.Source, &jpeg.Options{Quality: i.Quality})
		} else if i.SourceRGBA != nil {
			jpeg.Encode(newImg, i.SourceRGBA, &jpeg.Options{Quality: i.Quality})
		}
	}

	defer newImg.Close()
}

func ReadImage(fullPath string, altCheck bool) (*ImageReader, error) {
	r, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open: %s", err)
	}

	fi, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %s", err)
	}

	defer r.Close()
	fn := r.Name()
	i := &ImageReader{
		FileName: r.Name(),
		Format:   GetImageExt(fn),
		FileSize: fi.Size(),
	}

	var img image.Image
	if i.Format == "jpg" || i.Format == "jpeg" {
		img, err = jpeg.Decode(r)
		if err != nil {
			if altCheck == true {
				newFp := strings.ReplaceAll(fullPath, "jpeg", "png")
				newFp = strings.ReplaceAll(newFp, "jpg", "png")
				CopyFile(fullPath, newFp)
				return ReadImage(newFp, false)
			}

			return nil, fmt.Errorf("failed to decode: %s", err)
		}
	} else if i.Format == "png" {
		img, err = png.Decode(r)
		if err != nil {
			if altCheck == true {
				newFp := strings.ReplaceAll(fullPath, "png", "jpg")
				CopyFile(fullPath, newFp)
				return ReadImage(newFp, false)
			}

			return nil, fmt.Errorf("failed to decode: %s", err)
		}
	} else if i.Format == "gif" {
		nr := bufio.NewReader(r)
		g, err := gif.DecodeAll(nr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode GIF: %s", err)
		}
		if len(g.Image) == 0 {
			return nil, fmt.Errorf("failed to decode GIF: %s", err)
		}
		img = g.Image[0]
	}

	//Remove dot period
	if i.Format[:1] == "." {
		i.Format = i.Format[1:]
	}

	i.Image = &img
	i.Bounds = img.Bounds()
	i.Width = i.Bounds.Dx()
	i.Height = i.Bounds.Dy()

	return i, nil
}

func OverlayImage(baseImg *ImageReader, frameImage *ImageReader, fit *string, resizeOpts *ImageResizeOptions) *image.RGBA {
	if baseImg == nil || frameImage == nil {
		return nil
	}

	w := 1500
	h := 1500

	whiteBg := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	rect := image.Rect(0, 0, w, h)
	mergedImg := image.NewRGBA(rect)

	//1. Create an Empty Rectangle
	draw.Draw(mergedImg, mergedImg.Bounds(), &image.Uniform{whiteBg}, image.ZP, draw.Src)

	fi := *baseImg.Image
	//2. Forced base Image to have white background before resized
	//baseImgRect := image.Rect(0, 0, baseImg.Width, baseImg.Height)
	//baseImgWithBg := image.Image
	//draw.Draw(*baseImg.Image, baseImg.Bounds, &image.Uniform{whiteBg}, image.ZP, draw.Src)

	var fir image.Image
	if baseImg.Width == baseImg.Height || baseImg.Width > baseImg.Height {
		fir = resize.Resize(uint(w), 0, fi, resize.Lanczos3)
	} else {
		fir = resize.Resize(0, uint(h), fi, resize.Lanczos3)
	}

	fb := fir.Bounds()
	var offset image.Point
	if fb.Dx() > fb.Dy() {
		// rectangle
		offset.Y = (1500 - fb.Dy()) / 2
	} else if fb.Dy() > fb.Dx() {
		offset.X = (1500 - fb.Dx()) / 2
	}

	draw.Draw(mergedImg, fb.Add(offset), fir, image.ZP, draw.Src)

	si := *frameImage.Image
	var sir image.Image
	if fit != nil {
		ffit := enum.ImageFit(*fit)
		if ffit.IsValid() {
			switch ffit {
				case enum.ImageFitFill:
					sir = resize.Resize(1500, 1500, si, resize.Lanczos3)
				case enum.ImageFitContain:
					sir = resize.Resize(0, 1500, si, resize.Lanczos3)
				default:
					sir = resize.Resize(1500, 0, si, resize.Lanczos3)
			}
		} else {
			sir = resize.Resize(1500, 1500, si, resize.Lanczos3)
		}
	} else {
		sir = resize.Resize(1500, 1500, si, resize.Lanczos3)
	}
	//Overlay Fit

	sb := sir.Bounds()

	//Append second image over the first image
	draw.Draw(mergedImg, sb, sir, image.ZP, draw.Over)

	if resizeOpts != nil {
		if resizeOpts.MaintainRatio {
			if resizeOpts.Width > 0 {
				//Using width to resize with maintained ratio
				resized := resize.Resize(resizeOpts.Width, 0, mergedImg, resize.Lanczos3)
				if img, ok := resized.(*image.RGBA); ok {
					return img
				}
			} else if resizeOpts.Height > 0 {
				//Using height to resize with maintained ratio
				resized := resize.Resize(0, resizeOpts.Height, mergedImg, resize.Lanczos3)
				if img, ok := resized.(*image.RGBA); ok {
					return img
				}
			}
		} else {
			//Forced to resize using the actual width and height
			resized := resize.Resize(resizeOpts.Width, resizeOpts.Height, mergedImg, resize.Lanczos3)
			if img, ok := resized.(*image.RGBA); ok {
				return img
			}
		}
	}

	return mergedImg
}

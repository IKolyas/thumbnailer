package image

import (
	"crypto/sha256"
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

type Action string

const (
	ImageActionFill Action = "fill"
)

type Interface interface {
	Resize(width, height int) error
	Thumbnail(width, height int) error
	Export(format vips.ImageType) ([]byte, error)
	Fill()
	Convert()
}

type ImgData struct {
	ImageURL string
	Width    int
	Height   int
	Format   vips.ImageType
	Action   Action
}

func (img *ImgData) String() string {
	hash := sha256.New()
	hash.Write(fmt.Appendf(nil, "%s|%d|%d|%v|%s", img.ImageURL, img.Width, img.Height, img.Format, img.Action))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

type Image struct {
	VipsImg *vips.ImageRef
	Interface
}

func NewImage(imgData []byte) (*Image, error) {
	img, err := vips.NewImageFromBuffer(imgData)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	return &Image{VipsImg: img}, nil
}

func (i *Image) Fill(imgData *ImgData) ([]byte, error) {
	if err := i.resize(imgData.Width, imgData.Height); err != nil {
		return nil, fmt.Errorf("failed to resize image: %w", err)
	}

	if err := i.thumbnail(imgData.Width, imgData.Height); err != nil {
		return nil, fmt.Errorf("failed to thumbnail image: %w", err)
	}

	result, err := i.export()
	if err != nil {
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return result, nil
}

func (i *Image) resize(width, height int) error {
	scale := calculateScale(i.VipsImg, width, height)
	err := i.VipsImg.Resize(scale, vips.KernelLanczos3)
	if err != nil {
		return fmt.Errorf("failed to resize image: %w", err)
	}
	return nil
}

func (i *Image) thumbnail(width, height int) error {
	if width > 0 && height > 0 {
		err := i.VipsImg.Thumbnail(width, height, vips.InterestingCentre)
		if err != nil {
			return fmt.Errorf("failed to thumbnail image: %w", err)
		}
	}
	return nil
}

func (i *Image) export() ([]byte, error) {
	format := i.VipsImg.Metadata().Format

	switch format {
	case vips.ImageTypeJPEG:
		params := vips.NewJpegExportParams()
		params.Quality = 85
		params.OptimizeCoding = true
		imageBytes, _, err := i.VipsImg.ExportJpeg(params)
		if err != nil {
			return nil, fmt.Errorf("failed to export JPEG: %w", err)
		}
		return imageBytes, nil

	case vips.ImageTypePNG:
		params := vips.NewPngExportParams()
		params.Compression = 6
		params.Interlace = false
		imageBytes, _, err := i.VipsImg.ExportPng(params)
		if err != nil {
			return nil, fmt.Errorf("failed to export PNG: %w", err)
		}
		return imageBytes, nil

	case vips.ImageTypeWEBP:
		params := vips.NewWebpExportParams()
		params.Quality = 80
		params.Lossless = false
		params.ReductionEffort = 4
		imageBytes, _, err := i.VipsImg.ExportWebp(params)
		if err != nil {
			return nil, fmt.Errorf("failed to export WebP: %w", err)
		}
		return imageBytes, nil

	case vips.ImageTypeGIF, vips.ImageTypeTIFF, vips.ImageTypeBMP:
		// Для этих форматов используем стандартный экспорт
		imageBytes, _, err := i.VipsImg.ExportNative()
		if err != nil {
			return nil, fmt.Errorf("failed to export %v: %w", format, err)
		}
		return imageBytes, nil

	case vips.ImageTypeAVIF, vips.ImageTypeHEIF, vips.ImageTypeJP2K, vips.ImageTypeJXL:
		// Современные форматы с настройками по умолчанию
		imageBytes, _, err := i.VipsImg.ExportNative()
		if err != nil {
			return nil, fmt.Errorf("failed to export %v: %w", format, err)
		}
		return imageBytes, nil

	case vips.ImageTypePDF, vips.ImageTypeSVG, vips.ImageTypeMagick:
		// Векторные и специальные форматы
		return nil, fmt.Errorf("vector formats (%v) are not supported for export", format)

	case vips.ImageTypeUnknown:
		return nil, fmt.Errorf("unknown image format")

	default:
		// На случай добавления новых форматов в будущем
		return nil, fmt.Errorf("unsupported image format: %v", format)
	}
}

// вычисляет коэффициент масштабирования с сохранением пропорций.
func calculateScale(img *vips.ImageRef, width, height int) float64 {
	if width == 0 && height == 0 {
		return 1.0
	}

	switch {
	case width == 0 && height == 0:
		return 1.0
	case width > 0 && height > 0:
		scaleW := float64(width) / float64(img.Width())
		scaleH := float64(height) / float64(img.Height())
		return min(scaleW, scaleH)
	case width > 0:
		return float64(width) / float64(img.Width())
	default:
		return float64(height) / float64(img.Height())
	}
}

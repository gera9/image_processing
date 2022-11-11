package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
)

type Response struct {
	// Image contains the url to
	// the new processed image.
	Image string `json:"image"`
	Size  string `json:"size"`
	Type  string `json:"type"`
}

// EncodeToB64 encodes the image to base 64.
func EncodeToB64(data []byte, contentType string) (string, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("unable to decode jpeg: %w", err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return "", fmt.Errorf("unable to encode jpeg: %w", err)
	}

	data = buf.Bytes()

	imgBase64Str := base64.StdEncoding.EncodeToString(data)

	return imgBase64Str, nil
}

// CompressImage changes the mime type to jpeg and compresses the image.
func CompressImage(buffer []byte, quality int) ([]byte, error) {
	converted, err := bimg.NewImage(buffer).Convert(bimg.JPEG)
	if err != nil {
		return nil, err
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return processed, nil
}

// GenerateNewImageName generates an unique image name.
func GenerateNewImageName(extension string) string {
	return fmt.Sprintf("%s.%s", strings.Replace(uuid.New().String(), "-", "", -1), extension)
}

// GetImageExtension returns the image extension (jpeg or png).
func GetImageExtension(image []byte) string {
	return bimg.DetermineImageTypeName(image)
}

// WriteImage writes image to disk.
func WriteImage(img []byte, dst string) error {
	return os.WriteFile(dst, img, 0644)
}

func BuidlResponse(img []byte, imageName string) (*Response, error) {
	imgStats, err := os.Stat(fmt.Sprintf("pkg/uploads/%s", imageName))
	if err != nil {
		return nil, err
	}

	imgSize := imgStats.Size()

	return &Response{
		Image: fmt.Sprintf("http://localhost:3000/uploads/%s", imageName),
		Type:  GetImageExtension(img),
		Size:  fmt.Sprintf("%d", imgSize/1000),
	}, nil
}

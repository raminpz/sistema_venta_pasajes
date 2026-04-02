package pkg

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode genera un QR en base64 a partir de un string de datos.
func GenerateQRCode(data string, size int) (string, error) {
	if size <= 0 {
		size = 256
	}
	png, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		return "", fmt.Errorf("error generando QR: %w", err)
	}
	b64 := base64.StdEncoding.EncodeToString(png)
	return b64, nil
}

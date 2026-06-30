package dto

// GenerateQRResponse adalah respons dari endpoint generate QR
type GenerateQRResponseDto struct {
	QRString string `json:"qr_string"` // String yang akan diencode ke QR Code
}

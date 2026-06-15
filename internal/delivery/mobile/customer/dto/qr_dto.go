package dto

// GenerateQRResponse adalah respons dari endpoint generate QR
type GenerateQRResponse struct {
	QRString string `json:"qr_string"` // String yang akan diencode ke QR Code
}

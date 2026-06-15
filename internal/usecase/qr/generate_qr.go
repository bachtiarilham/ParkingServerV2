package qr

import (
	"context"
	"fmt"
	"time"
)

type GenerateQRInput struct {
	UserID int64 // Didapat dari context JWT
	// LocationID int64 // Jika ingin menyertakan lokasi, ambil dari parameter atau context
}

type GenerateQROutput struct {
	QRString string
}

type GenerateQRUseCase struct {
	// Tidak memerlukan repository untuk sementara, hanya logic generate string
	// Jika perlu validasi lokasi, tambahkan repo
}

func NewGenerateQRUseCase() *GenerateQRUseCase {
	return &GenerateQRUseCase{}
}

func (uc *GenerateQRUseCase) Execute(ctx context.Context, input GenerateQRInput) (GenerateQROutput, error) {
	// LINESPOT_QR:USR_123456789:LOC_A01:TS_1718234567
	qrString := fmt.Sprintf("LINESPOT_QR:USR_%d:TS_%d", input.UserID, time.Now().Unix())

	return GenerateQROutput{
		QRString: qrString,
	}, nil
}

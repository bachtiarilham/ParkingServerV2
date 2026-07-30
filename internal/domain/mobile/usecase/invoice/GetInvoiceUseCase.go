package invoice

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	core "modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/invoice"
	"modulegue/internal/domain/mobile/repository"
)

type GetInvoiceUseCase struct {
	invoiceRepo repository.InvoiceRepository
}

func NewGetInvoiceUseCase(
	invoiceRepo repository.InvoiceRepository,
) *GetInvoiceUseCase {
	return &GetInvoiceUseCase{
		invoiceRepo: invoiceRepo,
	}
}

func (uc *GetInvoiceUseCase) Execute(ctx context.Context, req string) (*model.UniversalInvoiceResponseModel, error) {
	result, err := uc.invoiceRepo.GetInvoice(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}

	doc, err := core.GenerateInvoicePDF(*result)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pdf: %w", err)
	}

	dirPath := "D:/parking_data/invoice"
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 6. Buat nama file PDF & Path Lengkap
	fileName := fmt.Sprintf("%v.pdf", result.TrxCode)
	fullPath := filepath.Join(dirPath, fileName)

	// 7. Simpan File PDF ke Local Disk
	if err := os.WriteFile(fullPath, doc, 0644); err != nil {
		return nil, fmt.Errorf("failed to save pdf to disk: %w", err)
	}
	result.InvoiceUrl = fileName

	return result, nil
}

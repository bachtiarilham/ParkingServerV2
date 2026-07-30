package laporan

// import (
// 	"context"
// 	"fmt"

// 	model "modulegue/internal/domain/mobile/model/laporan"
// 	"modulegue/internal/domain/mobile/repository"
// )

// type GetLaporanUseCase struct {
// 	laporanRepo repository.LaporanRepository
// }

// func NewGetLaporanUseCase(
// 	laporanRepo repository.LaporanRepository,
// ) *GetLaporanUseCase {
// 	return &GetLaporanUseCase{
// 		laporanRepo: laporanRepo,
// 	}
// }

// func (uc *GetLaporanUseCase) Execute(ctx context.Context, reqModel model.LaporanRequestModel) (*model.LaporanModel, error) {
// 	result, err := uc.laporanRepo.GetLaporan(ctx, reqModel)
// 	if err != nil {
// 		return nil, fmt.Errorf("get laporan: %w", err)
// 	}

// 	if result == nil {
// 		return nil, fmt.Errorf("get laporan: empty result")
// 	}

// 	return result, nil
// }

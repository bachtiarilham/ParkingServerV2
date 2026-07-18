package home

import (
	dto "modulegue/internal/data/website/remote/dto/home"
	model "modulegue/internal/domain/web/model/home"
)

func ToHomeResponseDto(src *model.HomeResponseModel) *dto.HomeResponseData {
	if src == nil {
		return nil
	}

	return &dto.HomeResponseData{
		NIK:             src.NIK,
		NoTelp:          src.NoTelp,
		GrafikPerDay:    toGrafikPerDayDtoSlice(src.GrafikPerDay),
		TransaksiAkhir:  toTransaksiAkhirDtoSlice(src.TransaksiAkhir),
		GrafikPerMinggu: toGrafikPerMingguDto(src.GrafikPerMinggu),
		ListProvinsi:    toProvinsiDtoSlice(src.ListProvinsi),
		ListShift:       toShiftDtoSlice(src.ListShift),
		Alamat:          src.Alamat,
		ListKecamatan:   toKecamatanDtoSlice(src.ListKecamatan),
		PetugasLapangan: toPetugasLapanganDtoSlice(src.PetugasLapangan),
		ID:              src.ID,
		StatsGlobal:     toStatsGlobalDto(src.StatsGlobal),
		Token:           src.Token,
		ListKabupaten:   toKabupatenDtoSlice(src.ListKabupaten),
		Nama:            src.Nama,
		Foto:            src.Foto,
		ListZona:        toZonaDtoSlice(src.ListZona),
		ListRole:        toRoleDtoSlice(src.ListRole),
		ListDesa:        toDesaDtoSlice(src.ListDesa),
	}
}

func ToHomeResponseModel(src *dto.HomeResponseData) *model.HomeResponseModel {
	if src == nil {
		return nil
	}

	return &model.HomeResponseModel{
		NIK:             src.NIK,
		NoTelp:          src.NoTelp,
		GrafikPerDay:    toGrafikPerDayModelSlice(src.GrafikPerDay),
		TransaksiAkhir:  toTransaksiAkhirModelSlice(src.TransaksiAkhir),
		GrafikPerMinggu: toGrafikPerMingguModel(src.GrafikPerMinggu),
		ListProvinsi:    toProvinsiModelSlice(src.ListProvinsi),
		ListShift:       toShiftModelSlice(src.ListShift),
		Alamat:          src.Alamat,
		ListKecamatan:   toKecamatanModelSlice(src.ListKecamatan),
		PetugasLapangan: toPetugasLapanganModelSlice(src.PetugasLapangan),
		ID:              src.ID,
		StatsGlobal:     toStatsGlobalModel(src.StatsGlobal),
		Token:           src.Token,
		ListKabupaten:   toKabupatenModelSlice(src.ListKabupaten),
		Nama:            src.Nama,
		Foto:            src.Foto,
		ListZona:        toZonaModelSlice(src.ListZona),
		ListRole:        toRoleModelSlice(src.ListRole),
		ListDesa:        toDesaModelSlice(src.ListDesa),
	}
}

func toGrafikPerDayDtoSlice(src *[]model.GrafikPerDayModel) *[]dto.GrafikPerDayDto {
	if src == nil {
		return nil
	}
	items := make([]dto.GrafikPerDayDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.GrafikPerDayDto(item))
	}
	return &items
}

func toGrafikPerDayModelSlice(src *[]dto.GrafikPerDayDto) *[]model.GrafikPerDayModel {
	if src == nil {
		return nil
	}
	items := make([]model.GrafikPerDayModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.GrafikPerDayModel(item))
	}
	return &items
}

func toTransaksiAkhirDtoSlice(src *[]model.TransaksiAkhirModel) *[]dto.TransaksiAkhirDto {
	if src == nil {
		return nil
	}
	items := make([]dto.TransaksiAkhirDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.TransaksiAkhirDto(item))
	}
	return &items
}

func toTransaksiAkhirModelSlice(src *[]dto.TransaksiAkhirDto) *[]model.TransaksiAkhirModel {
	if src == nil {
		return nil
	}
	items := make([]model.TransaksiAkhirModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.TransaksiAkhirModel(item))
	}
	return &items
}

func toGrafikPerMingguDto(src *model.GrafikPerMingguModel) *dto.GrafikPerMingguDto {
	if src == nil {
		return nil
	}
	return &dto.GrafikPerMingguDto{
		Motor: src.Motor,
		Mobil: src.Mobil,
		Total: src.Total,
	}
}

func toGrafikPerMingguModel(src *dto.GrafikPerMingguDto) *model.GrafikPerMingguModel {
	if src == nil {
		return nil
	}
	return &model.GrafikPerMingguModel{
		Motor: src.Motor,
		Mobil: src.Mobil,
		Total: src.Total,
	}
}

func toProvinsiDtoSlice(src *[]model.ProvinsiModel) *[]dto.ProvinsiDto {
	if src == nil {
		return nil
	}
	items := make([]dto.ProvinsiDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.ProvinsiDto(item))
	}
	return &items
}

func toProvinsiModelSlice(src *[]dto.ProvinsiDto) *[]model.ProvinsiModel {
	if src == nil {
		return nil
	}
	items := make([]model.ProvinsiModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.ProvinsiModel(item))
	}
	return &items
}

func toShiftDtoSlice(src *[]model.ShiftModel) *[]dto.ShiftDto {
	if src == nil {
		return nil
	}
	items := make([]dto.ShiftDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.ShiftDto(item))
	}
	return &items
}

func toShiftModelSlice(src *[]dto.ShiftDto) *[]model.ShiftModel {
	if src == nil {
		return nil
	}
	items := make([]model.ShiftModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.ShiftModel(item))
	}
	return &items
}

func toKecamatanDtoSlice(src *[]model.KecamatanModel) *[]dto.KecamatanDto {
	if src == nil {
		return nil
	}
	items := make([]dto.KecamatanDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.KecamatanDto(item))
	}
	return &items
}

func toKecamatanModelSlice(src *[]dto.KecamatanDto) *[]model.KecamatanModel {
	if src == nil {
		return nil
	}
	items := make([]model.KecamatanModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.KecamatanModel(item))
	}
	return &items
}

func toPetugasLapanganDtoSlice(src *[]model.PetugasLapanganModel) *[]dto.PetugasLapanganDto {
	if src == nil {
		return nil
	}
	items := make([]dto.PetugasLapanganDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.PetugasLapanganDto(item))
	}
	return &items
}

func toPetugasLapanganModelSlice(src *[]dto.PetugasLapanganDto) *[]model.PetugasLapanganModel {
	if src == nil {
		return nil
	}
	items := make([]model.PetugasLapanganModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.PetugasLapanganModel(item))
	}
	return &items
}

func toStatsGlobalDto(src *model.StatsGlobalModel) *dto.StatsGlobalDto {
	if src == nil {
		return nil
	}
	return &dto.StatsGlobalDto{
		JmlTransaksi: src.JmlTransaksi,
		TotalPAD:     src.TotalPAD,
		JmlPetugas:   src.JmlPetugas,
		Saldo:        src.Saldo,
	}
}

func toStatsGlobalModel(src *dto.StatsGlobalDto) *model.StatsGlobalModel {
	if src == nil {
		return nil
	}
	return &model.StatsGlobalModel{
		JmlTransaksi: src.JmlTransaksi,
		TotalPAD:     src.TotalPAD,
		JmlPetugas:   src.JmlPetugas,
		Saldo:        src.Saldo,
	}
}

func toKabupatenDtoSlice(src *[]model.KabupatenModel) *[]dto.KabupatenDto {
	if src == nil {
		return nil
	}
	items := make([]dto.KabupatenDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.KabupatenDto(item))
	}
	return &items
}

func toKabupatenModelSlice(src *[]dto.KabupatenDto) *[]model.KabupatenModel {
	if src == nil {
		return nil
	}
	items := make([]model.KabupatenModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.KabupatenModel(item))
	}
	return &items
}

func toZonaDtoSlice(src *[]model.ZonaModel) *[]dto.ZonaDto {
	if src == nil {
		return nil
	}
	items := make([]dto.ZonaDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.ZonaDto(item))
	}
	return &items
}

func toZonaModelSlice(src *[]dto.ZonaDto) *[]model.ZonaModel {
	if src == nil {
		return nil
	}
	items := make([]model.ZonaModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.ZonaModel(item))
	}
	return &items
}

func toRoleDtoSlice(src *[]model.RoleModel) *[]dto.RoleDto {
	if src == nil {
		return nil
	}
	items := make([]dto.RoleDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.RoleDto(item))
	}
	return &items
}

func toRoleModelSlice(src *[]dto.RoleDto) *[]model.RoleModel {
	if src == nil {
		return nil
	}
	items := make([]model.RoleModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.RoleModel(item))
	}
	return &items
}

func toDesaDtoSlice(src *[]model.DesaModel) *[]dto.DesaDto {
	if src == nil {
		return nil
	}
	items := make([]dto.DesaDto, 0, len(*src))
	for _, item := range *src {
		items = append(items, dto.DesaDto(item))
	}
	return &items
}

func toDesaModelSlice(src *[]dto.DesaDto) *[]model.DesaModel {
	if src == nil {
		return nil
	}
	items := make([]model.DesaModel, 0, len(*src))
	for _, item := range *src {
		items = append(items, model.DesaModel(item))
	}
	return &items
}

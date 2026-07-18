package setting

import (
	dto "modulegue/internal/data/website/remote/dto/setting"
	model "modulegue/internal/domain/web/model/setting"
)

func ToAddParlokRequestDto(src *model.AddParlokRequestModel) *dto.AddParlokRequestDto {
	if src == nil {
		return nil
	}
	return &dto.AddParlokRequestDto{
		NamaParlok:   src.NamaParlok,
		JalanParlok:  src.JalanParlok,
		IDZona:       src.IDZona,
		IDArea:       src.IDArea,
		IDDes:        src.IDDes,
		IDKec:        src.IDKec,
		IDKab:        src.IDKab,
		IDProv:       src.IDProv,
		LatMinArea:   src.LatMinArea,
		LatMaxArea:   src.LatMaxArea,
		LngMinArea:   src.LngMinArea,
		LngMaxArea:   src.LngMaxArea,
		AltitudeArea: src.AltitudeArea,
		CenterAreaX:  src.CenterAreaX,
		CenterAreaY:  src.CenterAreaY,
	}
}

func ToAddParlokRequestModel(src *dto.AddParlokRequestDto) *model.AddParlokRequestModel {
	if src == nil {
		return nil
	}
	return &model.AddParlokRequestModel{
		NamaParlok:   src.NamaParlok,
		JalanParlok:  src.JalanParlok,
		IDZona:       src.IDZona,
		IDArea:       src.IDArea,
		IDDes:        src.IDDes,
		IDKec:        src.IDKec,
		IDKab:        src.IDKab,
		IDProv:       src.IDProv,
		LatMinArea:   src.LatMinArea,
		LatMaxArea:   src.LatMaxArea,
		LngMinArea:   src.LngMinArea,
		LngMaxArea:   src.LngMaxArea,
		AltitudeArea: src.AltitudeArea,
		CenterAreaX:  src.CenterAreaX,
		CenterAreaY:  src.CenterAreaY,
	}
}

func ToRegisterRequestDto(src *model.RegisterRequestModel) *dto.RegisterRequest {
	if src == nil {
		return nil
	}
	return &dto.RegisterRequest{
		NIK:      src.NIK,
		Nama:     src.Nama,
		NoTelp:   src.NoTelp,
		Email:    src.Email,
		Username: src.Username,
		Password: src.Password,
		IDRole:   src.IDRole,
		Alamat:   src.Alamat,
		Foto:     src.Foto,
	}
}

func ToRegisterRequestModel(src *dto.RegisterRequest) *model.RegisterRequestModel {
	if src == nil {
		return nil
	}
	return &model.RegisterRequestModel{
		NIK:      src.NIK,
		Nama:     src.Nama,
		NoTelp:   src.NoTelp,
		Email:    src.Email,
		Username: src.Username,
		Password: src.Password,
		IDRole:   src.IDRole,
		Alamat:   src.Alamat,
		Foto:     src.Foto,
	}
}

func ToSaveScheduleRequestDto(src *model.SaveScheduleRequestModel) *dto.SaveScheduleRequestDto {
	if src == nil {
		return nil
	}
	return &dto.SaveScheduleRequestDto{
		ID:        src.ID,
		IDUser:    src.IDUser,
		IDLokasi:  src.IDLokasi,
		IDZona:    src.IDZona,
		IDArea:    src.IDArea,
		IDShift:   src.IDShift,
		DateAwal:  src.DateAwal,
		DateAkhir: src.DateAkhir,
	}
}

func ToSaveScheduleRequestModel(src *dto.SaveScheduleRequestDto) *model.SaveScheduleRequestModel {
	if src == nil {
		return nil
	}
	return &model.SaveScheduleRequestModel{
		ID:        src.ID,
		IDUser:    src.IDUser,
		IDLokasi:  src.IDLokasi,
		IDZona:    src.IDZona,
		IDArea:    src.IDArea,
		IDShift:   src.IDShift,
		DateAwal:  src.DateAwal,
		DateAkhir: src.DateAkhir,
	}
}

func ToSaveTarifRequestDto(src *model.SaveTarifRequestModel) *dto.SaveTarifRequestDto {
	if src == nil {
		return nil
	}
	return &dto.SaveTarifRequestDto{
		Tarif:           src.Tarif,
		KeteranganTarif: src.KeteranganTarif,
		IDLokasi:        src.IDLokasi,
	}
}

func ToSaveTarifRequestModel(src *dto.SaveTarifRequestDto) *model.SaveTarifRequestModel {
	if src == nil {
		return nil
	}
	return &model.SaveTarifRequestModel{
		Tarif:           src.Tarif,
		KeteranganTarif: src.KeteranganTarif,
		IDLokasi:        src.IDLokasi,
	}
}

func ToShowSelectedJukirRequestDto(src *model.ShowSelectedJukirRequestModel) *dto.ShowSelectedJukirRequestDto {
	if src == nil {
		return nil
	}
	return &dto.ShowSelectedJukirRequestDto{ID: src.ID}
}

func ToShowSelectedJukirRequestModel(src *dto.ShowSelectedJukirRequestDto) *model.ShowSelectedJukirRequestModel {
	if src == nil {
		return nil
	}
	return &model.ShowSelectedJukirRequestModel{ID: src.ID}
}

func ToShowSelectedJukirResponseDto(src *model.ShowSelectedJukirResponseModel) *dto.ShowSelectedJukirResponseDto {
	if src == nil {
		return nil
	}
	return &dto.ShowSelectedJukirResponseDto{
		NIK:      src.NIK,
		Username: src.Username,
		NoTelp:   src.NoTelp,
		SaldoMin: src.SaldoMin,
		IDStatus: src.IDStatus,
		IDRole:   src.IDRole,
		Alamat:   src.Alamat,
		ID:       src.ID,
		Email:    src.Email,
		Password: src.Password,
		Saldo:    src.Saldo,
		Nama:     src.Nama,
		Foto:     src.Foto,
	}
}

func ToShowSelectedJukirResponseModel(src *dto.ShowSelectedJukirResponseDto) *model.ShowSelectedJukirResponseModel {
	if src == nil {
		return nil
	}
	return &model.ShowSelectedJukirResponseModel{
		NIK:      src.NIK,
		Username: src.Username,
		NoTelp:   src.NoTelp,
		SaldoMin: src.SaldoMin,
		IDStatus: src.IDStatus,
		IDRole:   src.IDRole,
		Alamat:   src.Alamat,
		ID:       src.ID,
		Email:    src.Email,
		Password: src.Password,
		Saldo:    src.Saldo,
		Nama:     src.Nama,
		Foto:     src.Foto,
	}
}

func ToUpdateProfilRequestDto(src *model.UpdateProfilRequestModel) *dto.UpdateProfilRequestDto {
	if src == nil {
		return nil
	}
	return &dto.UpdateProfilRequestDto{
		IDJukir:  src.IDJukir,
		Nama:     cloneStringPtr(src.Nama),
		Username: cloneStringPtr(src.Username),
		Alamat:   cloneStringPtr(src.Alamat),
		NoTelp:   cloneStringPtr(src.NoTelp),
		Password: cloneStringPtr(src.Password),
		IDStatus: cloneStringPtr(src.IDStatus),
		Foto:     cloneStringPtr(src.Foto),
	}
}

func ToUpdateProfilRequestModel(src *dto.UpdateProfilRequestDto) *model.UpdateProfilRequestModel {
	if src == nil {
		return nil
	}
	return &model.UpdateProfilRequestModel{
		IDJukir:  src.IDJukir,
		Nama:     cloneStringPtr(src.Nama),
		Username: cloneStringPtr(src.Username),
		Alamat:   cloneStringPtr(src.Alamat),
		NoTelp:   cloneStringPtr(src.NoTelp),
		Password: cloneStringPtr(src.Password),
		IDStatus: cloneStringPtr(src.IDStatus),
		Foto:     cloneStringPtr(src.Foto),
	}
}

func cloneStringPtr(src *string) *string {
	if src == nil {
		return nil
	}
	value := *src
	return &value
}

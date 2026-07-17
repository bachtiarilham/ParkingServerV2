package topup

import (
	"strconv"
	"strings"

	dto "modulegue/internal/data/mobile/remote/dto/topup"
	model "modulegue/internal/domain/mobile/model/topup"
)

func ToTopupCreateRequestDto(src *model.TopupCreateRequestModel) *dto.TopupCreateRequestDto {
	if src == nil {
		return nil
	}

	return &dto.TopupCreateRequestDto{
		Amount:            src.Amount,
		PaymentMethodCode: src.PaymentMethodCode,
	}
}

func ToTopupCreateRequestModel(src *dto.TopupCreateRequestDto) *model.TopupCreateRequestModel {
	if src == nil {
		return nil
	}

	return &model.TopupCreateRequestModel{
		Amount:            src.Amount,
		PaymentMethodCode: src.PaymentMethodCode,
	}
}

func ToTopupCreateResponseDto(src *model.TopupCreateResponseModel) *dto.TopupCreateResponseDto {
	if src == nil {
		return nil
	}

	return &dto.TopupCreateResponseDto{
		TopupTransactionID: src.TopupTransactionID,
		TopupCode:          src.TopupCode,
		Amount:             src.Amount,
		AdminFee:           src.AdminFee,
		TotalAmount:        src.TotalAmount,
		PaymentMethodCode:  src.PaymentMethodCode,
		PaymentMethodName:  src.PaymentMethodName,
		PaymentStatusCode:  src.PaymentStatusCode,
		PaymentStatusName:  src.PaymentStatusName,
		QRString:           src.QRString,
		ExpiredAt:          src.ExpiredAt,
		CreatedAt:          src.CreatedAt,
	}
}

func ToTopupCreateResponseModel(src *dto.TopupCreateResponseDto) *model.TopupCreateResponseModel {
	if src == nil {
		return nil
	}

	return &model.TopupCreateResponseModel{
		TopupTransactionID: src.TopupTransactionID,
		TopupCode:          src.TopupCode,
		Amount:             src.Amount,
		AdminFee:           src.AdminFee,
		TotalAmount:        src.TotalAmount,
		PaymentMethodCode:  src.PaymentMethodCode,
		PaymentMethodName:  src.PaymentMethodName,
		PaymentStatusCode:  src.PaymentStatusCode,
		PaymentStatusName:  src.PaymentStatusName,
		QRString:           src.QRString,
		ExpiredAt:          src.ExpiredAt,
		CreatedAt:          src.CreatedAt,
	}
}

// func ToTopupStatusRequestDto(src *model.TopupStatusRequestModel) *dto.TopupStatusRequestDto {
// 	if src == nil {
// 		return nil
// 	}

// 	return &dto.TopupStatusRequestDto{
// 		TopupCode: src.TopupCode,
// 	}
// }

// func ToTopupStatusRequestModel(src *dto.TopupStatusRequestDto) *model.TopupStatusRequestModel {
// 	if src == nil {
// 		return nil
// 	}

// 	return &model.TopupStatusRequestModel{
// 		TopupCode: src.TopupCode,
// 	}
// }

func ToTopupStatusResponseDto(src *model.TopupStatusResponseModel) *dto.TopupStatusResponseDto {
	if src == nil {
		return nil
	}

	return &dto.TopupStatusResponseDto{
		TopupTransactionID: src.TopupTransactionID,
		TopupCode:          src.TopupCode,
		Amount:             src.Amount,
		AdminFee:           src.AdminFee,
		TotalAmount:        src.TotalAmount,
		PaymentMethodCode:  src.PaymentMethodCode,
		PaymentStatusCode:  src.PaymentStatusCode,
		PaymentMethodName:  src.PaymentMethodName,
		QRString:           src.QRString,
		PaidAt:             src.PaidAt,
		ExpiredAt:          src.ExpiredAt,
		FailedReason:       src.FailedReason,
		CompletedAt:        src.CompletedAt,
		CurrentBalance:     src.CurrentBalance,
	}
}

func ToTopupStatusResponseModel(src *dto.TopupStatusResponseDto) *model.TopupStatusResponseModel {
	if src == nil {
		return nil
	}

	return &model.TopupStatusResponseModel{
		TopupTransactionID: src.TopupTransactionID,
		TopupCode:          src.TopupCode,
		Amount:             src.Amount,
		AdminFee:           src.AdminFee,
		TotalAmount:        src.TotalAmount,
		PaymentMethodCode:  src.PaymentMethodCode,
		PaymentStatusCode:  src.PaymentStatusCode,
		PaymentMethodName:  src.PaymentMethodName,
		QRString:           src.QRString,
		PaidAt:             src.PaidAt,
		ExpiredAt:          src.ExpiredAt,
		FailedReason:       src.FailedReason,
		CompletedAt:        src.CompletedAt,
		CurrentBalance:     src.CurrentBalance,
	}
}

func ToQrisCallbackRequestDto(src *model.QrisCallbackRequestModel) *dto.QrisCallbackRequestDto {
	if src == nil {
		return nil
	}

	return &dto.QrisCallbackRequestDto{
		TransactionTime:   src.TransactionTime,
		TransactionStatus: src.TransactionStatus,
		TransactionID:     src.TransactionID,
		StatusMessage:     src.StatusMessage,
		StatusCode:        src.StatusCode,
		SignatureKey:      src.SignatureKey,
		PaymentType:       src.PaymentType,
		OrderID:           src.TopupCode,
		MerchantID:        src.MerchantID,
		GrossAmount:       src.GrossAmount,
		FraudStatus:       src.FraudStatus,
		Currency:          src.Currency,
		Acquirer:          src.Acquirer,
		SettlementTime:    src.SettlementTime,
	}
}

func ToQrisCallbackRequestModel(src *dto.QrisCallbackRequestDto) *model.QrisCallbackRequestModel {
	if src == nil {
		return nil
	}

	failedReason := ""
	switch strings.ToLower(src.TransactionStatus) {
	case "capture", "settlement", "pending":
		failedReason = ""
	default:
		if src.StatusMessage != "" {
			failedReason = src.StatusMessage
		} else {
			failedReason = src.TransactionStatus
		}
	}

	return &model.QrisCallbackRequestModel{
		TransactionTime:   src.TransactionTime,
		TransactionStatus: src.TransactionStatus,
		TransactionID:     src.TransactionID,
		StatusMessage:     src.StatusMessage,
		StatusCode:        src.StatusCode,
		SignatureKey:      src.SignatureKey,
		PaymentType:       src.PaymentType,
		TopupCode:         src.OrderID,
		MerchantID:        src.MerchantID,
		GrossAmount:       src.GrossAmount,
		FraudStatus:       src.FraudStatus,
		Currency:          src.Currency,
		Acquirer:          src.Acquirer,
		SettlementTime:    src.SettlementTime,
		FailedReason:      failedReason,
	}
}

func ToQrisCallbackResponseDto(src *model.QrisCallbackResponseModel) *dto.QrisCallbackResponseDto {
	if src == nil {
		return nil
	}

	return &dto.QrisCallbackResponseDto{
		TransactionID:     strconv.FormatInt(src.TopupTransactionID, 10),
		TransactionStatus: src.PaymentStatusCode,
		StatusMessage:     src.PaymentStatusCode,
		StatusCode:        "",
		PaymentType:       "qris",
		GrossAmount:       strconv.FormatInt(src.Amount, 10),
	}
}

func ToQrisCallbackResponseModel(src *dto.QrisCallbackResponseDto) *model.QrisCallbackResponseModel {
	if src == nil {
		return nil
	}

	amount, _ := strconv.ParseInt(src.GrossAmount, 10, 64)
	topupTransactionID, _ := strconv.ParseInt(src.TransactionID, 10, 64)

	return &model.QrisCallbackResponseModel{
		TopupTransactionID: topupTransactionID,
		Amount:             amount,
		PaymentStatusCode:  src.TransactionStatus,
	}
}

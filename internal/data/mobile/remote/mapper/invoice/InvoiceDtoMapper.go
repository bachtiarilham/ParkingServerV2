package invoice

import (
	dto "modulegue/internal/data/mobile/remote/dto/invoice"
	model "modulegue/internal/domain/mobile/model/invoice"
)

func ToUniversalInvoiceDto(src *model.UniversalInvoiceResponseModel) *dto.UniversalInvoiceResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.UniversalInvoiceResponseDto{
		TrxCode:         src.TrxCode,
		TransactionType: src.TransactionType,
		Title:           src.Title,
		Status:          src.Status,
		StatusText:      src.StatusText,
		CreatedAt:       src.CreatedAt,
		TotalAmount:     src.TotalAmount,
		InvoiceUrl:      src.InvoiceUrl,
		PriceBreakdown: dto.PriceBreakdownDto{
			BasePrice:      src.PriceBreakdown.BasePrice,
			AdminFee:       src.PriceBreakdown.AdminFee,
			TaxAmount:      src.PriceBreakdown.TaxAmount,
			DiscountAmount: src.PriceBreakdown.DiscountAmount,
			DiscountCode:   src.PriceBreakdown.DiscountCode,
			FinalTotal:     src.PriceBreakdown.FinalTotal,
		},
		PaymentMethod: dto.PaymentMethodDto{
			ChannelCode: src.PaymentMethod.ChannelCode,
			ChannelName: src.PaymentMethod.ChannelName,
			AccountNo:   src.PaymentMethod.AccountNo,
			IconURL:     src.PaymentMethod.IconURL,
		},
		CustomerInfo: dto.CustomerInfoDto{
			UserID:   src.CustomerInfo.UserID,
			FullName: src.CustomerInfo.FullName,
			Email:    src.CustomerInfo.Email,
			Phone:    src.CustomerInfo.Phone,
		},
	}

	if src.ParkingDetails != nil {
		out.ParkingDetails = &dto.ParkingInvoiceDetailDto{
			LocationName:  src.ParkingDetails.LocationName,
			LicensePlate:  src.ParkingDetails.LicensePlate,
			VehicleType:   src.ParkingDetails.VehicleType,
			CheckInTime:   src.ParkingDetails.CheckInTime,
			CheckOutTime:  src.ParkingDetails.CheckOutTime,
			DurationText:  src.ParkingDetails.DurationText,
			AttendantName: src.ParkingDetails.AttendantName,
		}
	}

	if src.WalletDetails != nil {
		out.WalletDetails = &dto.WalletInvoiceDetailDto{
			SenderName:       src.WalletDetails.SenderName,
			SenderAccount:    src.WalletDetails.SenderAccount,
			RecipientName:    src.WalletDetails.RecipientName,
			RecipientAccount: src.WalletDetails.RecipientAccount,
			BankRefNo:        src.WalletDetails.BankRefNo,
		}
	}

	if src.MembershipDetails != nil {
		out.MembershipDetails = &dto.MembershipInvoiceDetailDto{
			PackageName:     src.MembershipDetails.PackageName,
			PeriodStart:     src.MembershipDetails.PeriodStart,
			PeriodEnd:       src.MembershipDetails.PeriodEnd,
			MaxVehicles:     src.MembershipDetails.MaxVehicles,
			IsAutoRenew:     src.MembershipDetails.IsAutoRenew,
			NextBillingDate: src.MembershipDetails.NextBillingDate,
		}
	}

	return out
}
